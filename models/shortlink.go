package models

import (
	"time"
)

type Shortlink struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	OriginalURL   string    `json:"original_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsDeleted     bool      `json:"is_deleted"`
	IsDisabled    bool      `json:"is_disabled"`
	TotalVisits   int       `json:"total_visits"`
	TotalDuration int       `json:"total_duration"`

	InduceTitle    string `json:"induce_title"`
	InduceSubtitle string `json:"induce_subtitle"`
	InduceImageURL string `json:"induce_image_url"`
	InduceTemplate string `json:"induce_template"`
}

type ShortlinkCreateRequest struct {
	URL            string `json:"url" binding:"required,url"`
	InduceTitle    string `json:"induce_title"`
	InduceSubtitle string `json:"induce_subtitle"`
	InduceImageURL string `json:"induce_image_url"`
	InduceTemplate string `json:"induce_template"`
}

type ShortlinkResponse struct {
	ShortURL    string    `json:"short_url"`
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShortlinkListResponse struct {
	Items    []Shortlink `json:"items"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type Template struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Image    string `json:"image"`
	Icon     string `json:"icon"`
}

var TemplateLibrary = []Template{
	{
		ID:       "court",
		Name:     "法院传票",
		Title:    "您有一份法院传票待签收",
		Subtitle: "（2026粤0305民初xxx号）",
		Image:    "",
		Icon:     "⚖️",
	},
	{
		ID:       "express",
		Name:     "快递异常",
		Title:    "您的快递因收件地址不详被退回",
		Subtitle: "点击查看详情并处理",
		Image:    "",
		Icon:     "📦",
	},
	{
		ID:       "bank",
		Name:     "银行风控",
		Title:    "您的银行账户存在异常交易",
		Subtitle: "请立即核实，逾期将冻结账户",
		Image:    "",
		Icon:     "🏦",
	},
	{
		ID:       "subsidy",
		Name:     "政府补贴",
		Title:    "您有一笔政府补贴待领取",
		Subtitle: "金额：5000元，点击立即领取",
		Image:    "",
		Icon:     "💰",
	},
	{
		ID:       "license",
		Name:     "营业执照",
		Title:    "您的营业执照年度报告尚未提交",
		Subtitle: "请尽快处理，逾期将被吊销",
		Image:    "",
		Icon:     "📄",
	},
	{
		ID:       "traffic",
		Name:     "违章通知",
		Title:    "您的机动车违章记录已更新",
		Subtitle: "扣6分，罚款200元，点击处理",
		Image:    "",
		Icon:     "🚗",
	},
	{
		ID:       "social",
		Name:     "社保异常",
		Title:    "您的社保账户存在异常操作",
		Subtitle: "请立即核查，否则将影响待遇",
		Image:    "",
		Icon:     "🛡️",
	},
	{
		ID:       "redemption",
		Name:     "红包到账",
		Title:    "您有一笔红包已到账",
		Subtitle: "点击领取（仅限今日）",
		Image:    "",
		Icon:     "🧧",
	},
	{
		ID:       "vote",
		Name:     "投票结果",
		Title:    "您参与的投票结果已出",
		Subtitle: "点击查看您支持的候选人排名",
		Image:    "",
		Icon:     "🗳️",
	},
	{
		ID:       "custom",
		Name:     "自定义",
		Title:    "",
		Subtitle: "",
		Image:    "",
		Icon:     "✏️",
	},
}
