package live

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLiveStudents struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLiveStudents(ctx context.Context, svcCtx *svc.ServiceContext) *GetLiveStudents {
	return &GetLiveStudents{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLiveStudents) GetLiveStudents(req *types.GetLiveStudentsReq) (resp *types.GetLiveStudentsResp, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 50
	}
	count, err := l.svcCtx.LiveUserModel.GetStudentsCnt(l.ctx, req.StreamId, req.StudentName)
	if err != nil {
		logz.Errorf(l.ctx, "GetLiveStudents GetStudentsCnt err:%v", err)
		return nil, err
	}
	students, err := l.svcCtx.LiveUserModel.GetStudents(l.ctx, req.StreamId, req.Page, req.Limit, req.StudentName)
	resp = &types.GetLiveStudentsResp{
		List: make([]types.LiveStudentInfo, 0, len(students)),
		Pagination: types.Pagination{
			Total: count,
			Page:  req.Page,
			Limit: req.Limit,
		},
	}
	chatRoom, err := l.svcCtx.ChatRoomModel.FindSignalRoomByStreamId(l.ctx, req.StreamId)
	if err != nil {
		logz.Errorf(l.ctx, "GetLiveStudents FindSignalRoomByStreamId err:%v", err)
		return nil, err
	}
	for _, student := range students {
		resp.List = append(resp.List, types.LiveStudentInfo{
			UserId:   student.UserId,
			ImUserId: student.ImUserId,
			UserName: student.UserName,
			IsMute:   student.IsMute,
			EndAt:    student.MuteEndAt,
			RoomId:   chatRoom.RoomId,
		})
	}
	return
}

// 校验参数
func (l *GetLiveStudents) checkParams(req *types.GetLiveStudentsReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 GetLiveStudents 中编写

	return nil
}
