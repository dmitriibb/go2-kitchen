package buffers

import "github.com/dmitriibb/go2-kitchen/model"

// TODO add channel here in order to nox mix imports

var NewOrderItems = make(chan *model.OrderItem, 100)
var ReadyOrderItems = make(chan *model.OrderItem, 100)
