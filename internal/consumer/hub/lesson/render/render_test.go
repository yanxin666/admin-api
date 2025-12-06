package render

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/config"
	knowledgeModel "muse-admin/internal/model/knowledge"
	"muse-admin/internal/svc"
	"reflect"
	"testing"
)

func Test_convertLessonQuestionTTS(t *testing.T) {
	ctx, svcCtx, _ := svc.InitTest()
	type args struct {
		ctx        context.Context
		cos        *tencent.Cos
		questionId int64
		data       string
		conf       config.MiniMax
	}
	tests := []struct {
		name string
		args args
		want *knowledgeModel.LessonPoint
	}{
		{
			name: "",
			args: args{
				ctx: ctx,
				cos: tencent.NewCos(ctx, tencent.CocConf{
					SecretId:  svcCtx.Config.Oss.SecretId,
					SecretKey: svcCtx.Config.Oss.SecretKey,
					Appid:     svcCtx.Config.Oss.Appid,
					Bucket:    svcCtx.Config.Oss.Bucket,
					Region:    svcCtx.Config.Oss.Region,
				}),
				questionId: 1,
				data:       "王鹏飞你好，请听这个语音可以吗？请看下题唐僧娶了猪八戒",
				conf:       svcCtx.Config.MiniMax,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := convertLessonQuestionTTS(tt.args.ctx, tt.args.cos, tt.args.questionId, tt.args.data, tt.args.conf, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertLessonQuestionTTS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedule_handleQuestionTTS(t *testing.T) {
	ctx, svcCtx, _ := svc.InitTest()
	type fields struct {
		svcCtx           *svc.ServiceContext
		parentLessonData *knowledgeModel.Lesson
		cos              *tencent.Cos
	}
	type args struct {
		ctx     context.Context
		session sqlx.Session
		tts     *knowledgeModel.LessonQuestionTts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				svcCtx:           svcCtx,
				parentLessonData: nil,
				cos: tencent.NewCos(ctx, tencent.CocConf{
					SecretId:  svcCtx.Config.Oss.SecretId,
					SecretKey: svcCtx.Config.Oss.SecretKey,
					Appid:     svcCtx.Config.Oss.Appid,
					Bucket:    svcCtx.Config.Oss.Bucket,
					Region:    svcCtx.Config.Oss.Region,
				}),
			},
			args: args{
				ctx:     ctx,
				session: nil,
				tts:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schedule{
				svcCtx:           tt.fields.svcCtx,
				parentLessonData: tt.fields.parentLessonData,
				cos:              tt.fields.cos,
			}

			err := s.svcCtx.MysqlConnHub.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
				// question, _ := s.svcCtx.HubLessonQuestionModel.FindAll(ctx)
				// for _, v := range question {
				//	if v.Ask == "" {
				//		continue
				//	}
				//
				//	tts, err := convertLessonQuestionTTS(ctx, s.cos, v.Id, v.Ask, s.svcCtx.Config.MiniMax)
				//	if err != nil {
				//		fmt.Println(err)
				//		return err
				//	}
				//
				//	err = s.handleQuestionTTS(tt.args.ctx, session, tts)
				//	if err != nil {
				//		return err
				//	}
				//	fmt.Println(err)
				// }
				return nil
			})

			fmt.Println(err)
		})
	}
}
