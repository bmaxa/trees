package tree

type ITree interface {
  Size()uint
  Clear()
  Find(Key)Iterator
  Weight() (l,r uint)
  Height() uint
  Insert(Item)(Value,bool)
  Delete(Key) bool
  Validate()bool
  Begin()Iterator
  RBegin()Iterator
  End()Iterator
}

type IMulTree interface {
  ITree
  InsertEqual(Item)Value
  FindEqual(Key)(begin,end Iterator)
  DeleteIter(Iterator)bool
}

type Tree struct {
  Root *Node
  Size_ uint
}

type Node struct {
  Item
  Parent,Left,Right *Node
}

type Item struct {
  Key
  Value
}

type Key interface {
  Less(Key)bool
}

type Value interface {}

func (this* Tree) Size()uint {
  return this.Size_
}

func (this* Tree) Clear() {
  this.Root = nil
  this.Size_ = 0
}

func (this* Tree) Find(k Key) (rc Iterator) {
  var n,tmp = this.Root,(*Node)(nil)
  for n != nil {
    ret := k.Less(n.Key)
    if ret {
      n = n.Left
    } else {
      tmp = n
      n = n.Right
    }
  }
  if tmp != nil && !tmp.Less(k) {
    rc = NewIter(tmp)
  }
  return
}

func (this* Tree) FindEqual(k Key) (begin,end Iterator) {
  end = this.Find(k)
  if end == this.End() {
    begin = end
    return
  }
  for begin = end;begin != this.End() && begin.Prev() != this.End() && !begin.Prev().Node().Less(k); {
    begin = begin.Prev()
  }
  end = end.Next()
  return
}


func (this* Tree) Begin()Iterator {
  rc := (*Node)(nil)
  for n:=this.Root;n!=nil;n=n.Left {
    rc = n
  }
  return NewIter(rc)
}

func (this* Tree) RBegin()Iterator {
  rc := (*Node)(nil)
  for n:=this.Root;n!=nil;n=n.Right {
    rc = n
  }
  return NewIter(rc)
}

func (*Tree) End()Iterator {
  return NewIter(nil)
}

func (this* Tree) Weight()(Left,Right uint) {
  if this.Root == nil { return 0,0 }
  Left = Weight(this.Root.Left)
  Right = Weight(this.Root.Right)
  return
}

func Weight(n *Node)uint {
  if n == nil { return 0 }
  return Weight(n.Left)+Weight(n.Right)+1
}

func (this* Tree) Height()(rc uint) {
  if this.Size_ == 0 { return 0 }
  height(this.Root,0,&rc)
  return
}

func height(n *Node,d uint,max *uint){
  if n == nil {
    if *max < d {
      *max = d
    }
    return
  }
  height(n.Left,d+1,max)
  height(n.Right,d+1,max)
}

func ToString(n* Node,prefix string,isTail bool,df func(*Node)string)string {
  var tmp = prefix
  if isTail {
    tmp = tmp + "└── "
  } else {
    tmp = tmp + "├── "
  }
  if n == nil {
    tmp = tmp + "nil"
    return tmp
  }
  tmp = tmp + df(n)

  var tmp1 string
  if isTail {
    tmp1 = prefix + "    "
  } else {
    tmp1 = prefix + "│   "
  }
  
  tmp = tmp + ToString(n.Left,tmp1,false,df) + "\n"
  tmp = tmp + ToString(n.Right,tmp1,true,df)  
  return tmp
}

type Iterator struct {
  data *Node
}

func NewIter(n *Node) Iterator {
  return Iterator{n}
}

func (this Iterator) Value()Item {
  return this.data.Item
}

func (this Iterator) Node()*Node{
  return this.data
}

func (this Iterator) Next()Iterator{
  if this.data.Right != nil {
    this.data = this.data.Right
    for this.data.Left != nil {
      this.data = this.data.Left
    }
  } else {
    tmp := this.data.Parent
    for tmp != nil && this.data == tmp.Right {
      this.data = tmp
      tmp = tmp.Parent
    }
    this.data = tmp
  }
  return this
}

func (this Iterator) Prev()Iterator{
  if this.data.Left != nil {
    this.data = this.data.Left
    for this.data.Right != nil {
      this.data = this.data.Right
    }
  } else {
    tmp := this.data.Parent
    for tmp != nil && this.data == tmp.Left {
      this.data = tmp
      tmp = tmp.Parent
    }
    this.data = tmp
  }
  return this
}
