package dtypes

type ReplyGuyPost struct {
	ID      int
	Content string
	Author  Author `json:"author"`
}

type ReplyGuyComment struct {
	ID      int
	Content string
	Author  Author `json:"author"`
}

type ReplyGuyRequest struct {
	Comment       ReplyGuyComment `json:"comment"`
	ParentComment ReplyGuyComment `json:"parentComment"`
	ParentPost    ReplyGuyPost    `json:"parentPost"`
	Model         string          `json:"model"`
}

type OllamaRequest struct {
	Stream bool   `json:"stream"`
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
