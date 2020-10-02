package database

import (
	"database/sql"
	"fmt"
)

// ReadQuestionList 使用 u：用户名(ID) 查询 QuestionListTab 表。
// 列出所有答题
func (q Question) ReadQuestionList(u string) (tab []*QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.creator_id = "%v" ORDER BY question_list.id DESC`,
		u,
	)

	rows, err := Class.DB.Query(sq)
	if err != nil {
		return
	}
	defer rows.Close()

	tab, err = joinQuestionList(rows)
	return

}

// ReadQuestion 使用 i：问题ID(ID) 查询 QuestionListTab 表。
// TODO 联表查询 定义问题结构体
// 答题信息
func (q Question) ReadQuestion(i uint32) (data *QuestionListTab, err error) {

	sq := fmt.Sprintf(
		`SELECT question_list.* FROM question_list WHERE question_list.id = %v`,
		i,
	)

	row, err := Class.DB.Query(sq)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		return
	}
	data = new(QuestionListTab)
	err = row.Scan(&data.ID, &data.Question, &data.CreatorID, &data.Market)
	if err != nil {
		return
	}

	return

}

// ReadQuestionMarket 查询 QuestionListTab 表。
// 答题市场
func (q Question) ReadQuestionMarket() (tab []*QuestionListTab, err error) {

	rows, err := Class.DB.Query(
		`SELECT question_list.* FROM question_list WHERE question_list.market = true ORDER BY question_list.id DESC`,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	tab, err = joinQuestionList(rows)
	return

}

// joinQuestionList 复用
func joinQuestionList(rows *sql.Rows) (tab []*QuestionListTab, err error) {

	var data []*QuestionListTab
	for rows.Next() {

		data0 := new(QuestionListTab)
		err = rows.Scan(
			&data0.ID, &data0.Question, &data0.CreatorID, &data0.Market,
		)
		if err != nil {
			return
		}

		data = append(data, data0)

	}

	tab = data
	return

}

// WriteQuestionList 写入 QuestionListTab 表。
// 新建答题
func (q Question) WriteQuestionList(tab *QuestionListTab) (err error) {

	i, err := Class.DB.Prepare(
		`INSERT INTO question_list (id, question, creator_id, market) VALUES (?, ?, ?, ?)`,
	)
	if err != nil {
		return
	}
	defer i.Close()

	_, err = i.Exec(nil, tab.Question, tab.CreatorID, tab.Market) // ID 自增无需输入
	if err != nil {
		return
	}

	return

}