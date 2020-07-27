package lightweb

type Param struct {
	Key,Val string
}

type Params  []Param

func (ps Params) ByName(name string) (string) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Val
		}
	}
	return ""
}

var handlersSet []HandlerChain

type node struct {
	key uint8
	tag string	//动态路由标识
	handlers HandlerChain
	child []*node
	leaf bool //是否有一个路由
}

type MethodHash map[string]*node

func (mh MethodHash) getRoot(method string) *node {
	if mh[method] != nil{
		return mh[method]
	}

	//如果没有，自动创建
	p := newNode('/')
	mh[method] = p
	return p
}

func (n *node) binary_search(key uint8) *node {
	//遇到动态路由
	if len(n.child) == 1 && n.child[0].key == ':' {
		return n.child[0]
	}

	l := 0
	r := len(n.child) - 1
	for l < r {
		mid := (l + r) >> 1
		if n.child[mid].key > key {
			r = mid - 1
		}else {
			l = mid
		}
	}

	//路径不存在
	if l >= len(n.child) || n.child[l].key != key {
		return nil
	}

	return n.child[l]
}

func newNode(key uint8) *node {
	return &node{
		key:    key,
		tag:	"",
		child:  make([]*node,0),
		leaf:   false,
	}
}

//读取字符串直到'/'，并将i停在'/'的前一位
func walk(i *int,path string) (res string) {
	res = ""
	if path[*i] == ':' {
		*i++
	}

	for ;*i < len(path); *i++ {
		if path[*i] == '/' {
			*i--
			return
		}
		res = res + string(path[*i])
	}
	//读到末尾，防止越界
	*i--
	return
}

func (mh *MethodHash) addRouter(method string,path string, handlers HandlerChain) bool {
	rt := mh.getRoot(method)

	//创建路径
	p := rt.insert(path,1)
	if p == nil {
		return false
	}

	//注册路由
	p.handlers = append(p.handlers, handlers...)
	p.leaf = true
	return  true
}


func (n *node) insert(path string,i int) *node {
	//创建路径
	ch := n.binary_search(path[i])
	flag := false
	if ch == nil {
		ch = newNode(path[i])
		flag = true
	}

	if path[i] == ':' {
		if ch == nil {
				return nil
			}
		ch.tag = walk(&i,path)
	}

	var p *node
	if i == len(path) - 1 {
		if !flag {
			p = nil
		} else {
			p = ch
		}
	} else {
		p = ch.insert(path,i+1)
	}


	//创建成功，连接新建节点
	if flag && p != nil {
		n.child = append(n.child, ch)
		pos := len(n.child) - 1

		//重新排序
		for pos > 0 && n.child[pos].key < n.child[pos-1].key{
			//swap
			n.child[pos],n.child[pos-1] = n.child[pos-1],n.child[pos]
			pos--
		}
	}
	//失败则ch自动回收且放弃对树的更改

	return p
}

func (mh MethodHash) getRouter(method,path string) (handlers HandlerChain,params Params,ok bool) {
	rt := mh.getRoot(method)
	handlers,params,ok = rt.getHandlers(path)
	return
}

func (n *node) getHandlers(path string) (handlers HandlerChain,params Params,ok bool) {
	p := n
	for i := 1;i < len(path); i++ {
		nxt := p.binary_search(path[i])
		if nxt == nil {
			return  nil,nil,false
		}
		p = nxt

		//动态路由
		if p.key == ':' {
			val := walk(&i,path)
			params = append(params, Param{
				Key: p.tag,
				Val: val,
			})
		}
	}

	handlers = p.handlers
	//不是叶子节点
	if !p.leaf {
		return nil,nil,false
	}
	ok = true
 	return
}