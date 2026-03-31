package inputs

type AddCommentInput struct {
	Comment   string `json:"comment" db:"comment"`
	TaskId    int    `json:"task_id" db:"task_id"`
	CreatorId int    `json:"-"`
}
