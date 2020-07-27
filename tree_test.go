package lightweb

import (
	"fmt"
	"testing"
)




func h(c *Context) {
	fmt.Println("ok")
}

func dfs(n *node) {
	fmt.Printf("%c",n.key)
	for _,ch := range n.child {
		if ch != nil {
			dfs(ch)
		}
	}
}

func TestAdd(t *testing.T) {
	mh := make(MethodHash,5)
	hc := make(HandlerChain,1)
	hc[0] = h
	mh.addRouter("Get","/test",hc)
	hh,_,_ := mh.getRoot("Get").getHandlers("/test")
	dfs(mh.getRoot("Get"))
	var ctx *Context
	fmt.Println(len(hh))
	if len(hh) > 0 && hh[0] != nil{
		hh[0](ctx)
	}

}