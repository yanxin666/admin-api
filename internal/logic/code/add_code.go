package code

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/model/code"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"regexp"
	"strconv"
	"time"
)

type AddCode struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCode(ctx context.Context, svcCtx *svc.ServiceContext) *AddCode {
	return &AddCode{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCode) AddCode(req *types.AddCodeReq) error {
	codes, err := NewGenerateCode(l.ctx, l.svcCtx).BatchRedemptionCodes(l.ctx, int(req.CodeNumber))
	if err != nil {
		return err
	}
	successCodeNumber, err := l.InsertBatchCode(l.ctx, codes, req)
	if err != nil {
		return err
	}
	logc.Infof(l.ctx, "生成兑换码:%d", successCodeNumber)
	return nil
}

func (l *AddCode) InsertBatchCode(ctx context.Context, codes []string, req *types.AddCodeReq) (int, error) {
	// 记录成功入库的兑换码数量
	successfulInserts := 0
	datetime, _ := GetStandardDatetime(req.ValidDate)
	for i := 0; i < len(codes); i += batchSize {
		end := i + batchSize
		if end > len(codes) {
			end = len(codes)
		}
		batchCodes := codes[i:end]
		redeemCodeData := make([]*code.RedeemCode, 0, len(batchCodes))

		for _, v := range batchCodes {
			if req.Batch == 0 {
				currentTime := time.Now()
				formattedTime := currentTime.Format("20060102")
				formattedTimeInt, _ := strconv.ParseInt(formattedTime, 10, 64)
				req.Batch = formattedTimeInt
			}
			redeemCodeData = append(redeemCodeData, &code.RedeemCode{
				Code:            v,
				Status:          1,
				Batch:           req.Batch,
				ValidDate:       datetime,
				Source:          2,                   // 来源 1.外部导入 2.内容生成
				BenefitsGroupId: req.BenefitsGroupId, // 权益组表ID
			})
		}

		// 重试机制
		for attempt := 0; attempt < maxRetries; attempt++ {
			logc.Infof(ctx, "第 %d 次尝试插入兑换码", attempt+1)

			err := l.svcCtx.RedeemCodeModel.Inserts(l.ctx, redeemCodeData)
			if err == nil {
				// 插入成功，累加成功插入的数量
				successfulInserts += len(redeemCodeData)
				logc.Infof(ctx, "第 %d 次插入成功，已成功插入 %d 条数据", attempt+1, successfulInserts)
				break // 插入成功，处理下一批
			}

			logx.Errorf("插入失败, 错误: %v", err)

			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				// 处理唯一索引冲突，生成新兑换码
				conflictCodes := extractConflictCodes(err)
				logc.Infof(ctx, "检测到冲突的兑换码: %v", conflictCodes)

				if len(conflictCodes) > 0 {
					newCodes, err := NewGenerateCode(l.ctx, l.svcCtx).BatchRedemptionCodes(l.ctx, len(conflictCodes))
					if err != nil {
						logz.Errorf(ctx, "索引错误，生成兑换码错误，不重试: %v", err)
						return successfulInserts, err
					}
					logc.Infof(ctx, "生成新的兑换码: %v", newCodes)
					// 更新冲突兑换码
					for i, conflictCode := range conflictCodes {
						for j := range redeemCodeData {
							if redeemCodeData[j].Code == conflictCode {
								redeemCodeData[j].Code = newCodes[i]
							}
						}
					}
				}
			} else {
				// 非唯一索引错误直接返回，不继续重试
				logz.Errorf(ctx, "非唯一索引错误，不进行重试: %v", err)
				return successfulInserts, err
			}

			// 达到最大重试次数后仍未成功
			if attempt == maxRetries-1 {
				logz.Errorf(ctx, "已达到最大重试次数，仍未能成功插入兑换码")
				return successfulInserts, errors.New("超过最大重试次数，仍未能成功插入兑换码")
			}
		}
	}
	return successfulInserts, nil
}

// extractConflictCodes 从 MySQL 错误信息中提取冲突的兑换码
func extractConflictCodes(err error) []string {
	errMsg := err.Error()

	// 使用正则表达式提取 "Duplicate entry" 后面的兑换码
	re := regexp.MustCompile(`Duplicate entry '([A-Z0-9]+)' for key`)
	matches := re.FindStringSubmatch(errMsg)

	if len(matches) > 1 {
		// 返回提取的兑换码
		return []string{matches[1]}
	}
	return nil
}

const StandardDatetime = "2006-01-02 15:04:05"

// GetStandardDatetime 获取标准时间并确保时区为 UTC+8
func GetStandardDatetime(dateString string) (time.Time, error) {
	// 定义一个时区
	loc, err := time.LoadLocation("Asia/Shanghai") // 选择正确的时区
	if err != nil {
		return time.Date(2099, 12, 31, 0, 0, 0, 0, loc), nil // 返回默认日期
	}
	// 将输入的字符串解析为时间并应用时区
	standardTime, err := time.ParseInLocation(StandardDatetime, dateString, loc)
	if err != nil {
		return time.Date(2099, 12, 31, 0, 0, 0, 0, loc), nil // 返回默认日期
	}
	return standardTime, nil
}
