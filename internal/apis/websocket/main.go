package websocket

import (
	"net/http"
	"strconv"

	class "github.com/ELQASASystem/app/internal/app"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// srv Websocket 服务。文档： https://godoc.org/github.com/gorilla/websocket
type srv struct {
	connPool map[uint32][]*websocket.Conn // connPool 已连接的客户端
	qch      chan *class.Question         // qch 问题管道
}

// New 新建 Websocket 服务
func New() chan *class.Question {

	var (
		qch = make(chan *class.Question, 10)
		s   = &srv{connPool: map[uint32][]*websocket.Conn{}, qch: qch}
	)

	go s.start()

	return qch
}

// start 启动 Websocket 服务
func (w *srv) start() {

	http.HandleFunc("/question", w.handleWS)
	go w.sendQuestion()

	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Error().Err(err).Msg("失败")
	}

}

func (w *srv) handleWS(writer http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // 允许跨域

	wsconn, err := upgrader.Upgrade(writer, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}
	defer wsconn.Close()

	for {

		_, msg, err := wsconn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			return
		}

		action := string(msg)
		log.Info().Str("消息", action).Msg("收到 Websocket 消息")

		if action == "keep heart" {
			continue
		}

		i, _ := strconv.ParseUint(action, 10, 64)

		w.addConn(uint32(i), wsconn)
		defer w.rmConn(uint32(i), wsconn)
		log.Info().Uint64("问题ID", i).Msg("成功添加WS问题监听")

	}

}

// addConn 使用 i：问题ID(ID) 新增一个连入的客户端
func (w *srv) addConn(i uint32, c *websocket.Conn) {

	if _, ok := w.connPool[i]; !ok {
		w.connPool[i] = []*websocket.Conn{c}
		return
	}

	w.connPool[i] = append(w.connPool[i], c)
}

// rmConn 使用 i：问题ID(ID) 移出一个连入的客户端
func (w *srv) rmConn(i uint32, conn *websocket.Conn) {

	conns := w.connPool[i]

	for k, wsconn := range conns {
		if wsconn == conn {
			w.connPool[i] = append(w.connPool[i][:k], w.connPool[i][k+1:]...)
		}
	}

	if len(conns) == 0 {
		delete(w.connPool, i)
	}

	log.Info().Msg("WS客户端下线")

}
