package public

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/config"
	"os"
	"time"
)

// DownloadFiled 下载文件到本地
func DownloadFiled(ctx context.Context, oss config.Oss, id int64, name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		logc.Errorf(ctx, "获取文件路径失败,Err:%v", err)
		return "", err
	}

	downloadPath := dir + "/download/"

	if !checkFileExists(downloadPath) {
		err = os.Mkdir(downloadPath, 0755)
		if err != nil {
			logc.Errorf(ctx, "创建文件路径失败,Err:%v", err)
			return "", err
		}
	}

	filed := cast.ToString(time.Now().UnixMilli()) + name
	filename := downloadPath + filed

	cos := tencent.NewCos(ctx, tencent.CocConf{
		SecretId:  oss.SecretId,
		SecretKey: oss.SecretKey,
		Appid:     oss.Appid,
		Bucket:    oss.Bucket,
		Region:    oss.Region,
	})

	isTrue, err := cos.DownLoadToLoad("RD/files/", name, filename)
	if err != nil {
		logc.Errorf(ctx, "创建获取cos文件失败,导入id:%d,导入名称:%s,错误原因%v", id, name, err)
		return "", err
	}
	if !isTrue {
		logc.Errorf(ctx, "所下载的文件对象不存在,导入id:%d,导入名称:%s,错误原因%v", id, name, err)
		return "", errors.New("所下载的文件对象不存在")
	}
	return filename, nil
}

// checkFileExists 检查文件路径是否存在
func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
