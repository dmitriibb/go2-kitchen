package buffers

import "github.com/dmitriibb/go2-kitchen/internal/model"

var NewOrderItems = make(chan *model.OrderItem, 100)
var ReadyOrderItems = make(chan *model.OrderItem, 100)
