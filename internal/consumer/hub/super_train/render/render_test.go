package render

import (
	"context"
	"encoding/json"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"testing"
)

func TestRender_Render(t *testing.T) {
	ctx, svcCtx, _ := svc.InitTest()
	type fields struct {
		svcCtx *svc.ServiceContext
	}
	type args struct {
		ctx context.Context
		f   *types.WritePPTData
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
				svcCtx: svcCtx,
			},
			args: args{
				ctx: ctx,
				f:   nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Render{
				svcCtx: tt.fields.svcCtx,
			}
			str := `{"source":4,"env_type":2,"account":"wangpengfei","status":3,"course_type":1,"course_no":"C01-01-013","course_name":"在测一遍-课程、章节加字段了","lesson_name":"课节名称","subject":1,"unit":"单元名","image":"https://static-jx-admin.zmexing.com/C01-01-013/29ccb0dd-8c23-4465-8721-831ecabba23c.png","intro":"单元名","course_extra":{"title":"首页课程标题","button_text":"首页按钮文字","teacher_name":"首页老师名字","teacher_video":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/eb330db65145403705383610406/v.f100240.m3u8","teacher_video_web":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/eb330db65145403705383610406/v.f100040.mp4","home_page_guide_video":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/eb330db65145403705383610406/v.f100240.m3u8","course_background":"课程背景描述","course_appellation":"课程称谓描述","through":"穿越名","article_tts":"作文选字数tts","secrets_task":"作文选字数tts","article_title":"作文题目","outline_tts_text":["写作大纲tts文本","写作大纲tts文本"],"after_article_tts":"生成作文后的tts","writing_skills":"写作技法 大模型使用","article_center_idea":"作文中心思想 大模型使用"},"chapter":[{"no":"CH01-01-013-01","name":"章节一_章节标题","title":"章节精简版标题","image":"https://static-jx-admin.zmexing.com/C01-01-013/6464e282-18d7-405a-8870-83e696eb8a58.png","intro":"章节介绍","type":1,"sequence":1,"guide_video":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/e96f2de05145403705383591273/v.f100240.m3u8","teacher_id":1,"chapter_extra":{"guide_video":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/e26217ea5145403705383289864/v.f100240.m3u8","guide_task":[{"video_url":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/9e61e5a15145403705382639182/v.f100240.m3u8","video_web_url":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/9e61e5a15145403705382639182/v.f100240.m3u8","content":"引导内容1","images":[{"url":"https://static-jx-admin.zmexing.com/admin/29baaeec-0e82-4050-96cd-e83d24f4d1ba.png","theme":"事件主题","title":"事件标题1","trigger":{"type":2,"payload":"{\"video_url\":\"\",\"video_web_url\":\"\",\"light_events\":[{\"name\":\"namenamename\",\"position\":{\"x\":1,\"y\":1},\"time\":1111}]}"},"time":12},{"url":"https://static-jx-admin.zmexing.com/admin/77756262-a359-4dba-a963-83ebf67416b3.png","theme":"事件主题2","title":"事件标题2","trigger":{"type":1,"payload":"{\"video_url\":\"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/df8e21185145403705383119017/v.f100240.m3u8\",\"video_web_url\":\"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/df8e21185145403705383119017/v.f100240.m3u8\",\"light_events\":null}"},"time":22}]},{"video_url":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/e971f0475145403705383596392/v.f100240.m3u8","video_web_url":"https://1500040952.vod-qcloud.com/43c9e764vodtranscq1500040952/e971f0475145403705383596392/v.f100240.m3u8","content":"引导内容2","images":[{"url":"https://static-jx-admin.zmexing.com/admin/6bcc076c-b13e-4a03-a0fb-2082dc311bd6.jpeg","theme":"事件主题2-1","title":"事件标题2-1","trigger":{"type":2,"payload":"{\"video_url\":\"\",\"video_web_url\":\"\",\"light_events\":[{\"name\":\"lightEventsname\",\"position\":{\"x\":1,\"y\":1},\"time\":112}]}"},"time":12}]}],"writing_techniques":"写作技法写作技法","composition_title":"作文标题作文标题","writing_goals":"写作目标写作目标","full_star_advice":"文章页面满星建议文章页面满星建议","excellent_sample":"优秀范文示例优秀范文示例优秀范文示例优秀范文示例优秀范文示例优秀范文示例","article_top_tip_text":"文章页面顶部文本文章页面顶部文本文章页面顶部文本","background_url":"https://static-jx-admin.zmexing.com/C01-01-013/88af425a-9872-4d9f-86ea-332e7e7cb5ad.png","guide_video_desc":"引导页视频文案引导页视频文案引导页视频文案引导页视频文案"},"task":[{"task_no":"T01-01-013-01-01","title":"任务一_1","subtitle":"子标题","description":"描述","progress_prefix":"进度前缀","type":4,"sequence":1,"sub_task":[{"sub_task_no":"ST01-01-013-01-01-01","mode":1,"type":1,"title":"子任务一_标题","description":"描述","image_url":"https://static-jx-admin.zmexing.com/C01-01-013/2842f26e-db83-4dae-a9d1-b9d9ae33f640.png","sequence":1,"sub_extra":{"question":"题目","exemplary_sentences":["优秀句子示例优秀句子示例优秀句子示例优秀句子示例","优秀句子示例优秀句子示例优秀句子示例优秀句子示例"],"answer_key":["答案关键字答案关键字答案关键字"],"goals":"写作要求写作要求写作要求","storyline":"写作要求写作要求写作要求","approach":"写作思路写作思路写作思路写作思路","grade_rules":[]}}]}]}]}`
			var f *types.SuperTrainData
			_ = json.Unmarshal([]byte(str), &f)
			if err := r.Render(tt.args.ctx, f); (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
