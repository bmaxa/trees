package splay

import (
  "fmt"
  "github.com/bmaxa/trees/tree"
  "unsafe"
)

type SPLAY struct{
  tree.Tree
}

type Node struct {
  tree.Node
}

func uc(n *Node) *tree.Node{
  rc := (*tree.Node)(unsafe.Pointer(n))
  return rc
}
     
func dc(n *tree.Node)*Node{
  rc := (*Node)(unsafe.Pointer(n))
  return rc
}
          
func New()*SPLAY{
  return &SPLAY{}
}

func NewNode(i tree.Item)*tree.Node {
  n := tree.Node{ Item:i }
  return uc(&Node{ Node:n })
}

func (t* SPLAY) Validate()bool {
  return validate(dc(t.Root))
}

func (t SPLAY) String() string {
  if t.Root == nil { return "Empty SPLAY" }
  return tree.ToString(t.Root,"",true,to_string)
}

func (this* SPLAY) Insert(i tree.Item)(tree.Value,bool) {
  if this.Root == nil {
    this.Root = NewNode(i)
    this.Size_++
    return this.Root.Value,true
  }
  var prev *tree.Node
  var rc *tree.Node
  var ret bool
  for n := this.Root ; n != nil; {
    prev = n
    if ret = i.Less(n.Key);ret {
      n = n.Left
    } else {
      rc = n
      n = n.Right
    }
  }

  if rc != nil && !rc.Less(i.Key) {
    this.splay(rc)
    return rc.Value,false
  }  

  n := NewNode(i)
  this.Size_++
  if !ret {
    prev.Right = n
    n.Parent = prev
  } else {
    prev.Left = n
    n.Parent = prev
  }
  this.splay(n)
  return n.Value,true
}

func (this* SPLAY) Delete(k tree.Key)bool {
  n := this.Find(k).Node()
  if n == nil {
    return false
  } 
  this.delete(n)
  return true
}

func (this* SPLAY) delete(n *tree.Node) {
  if n.Left != nil && n.Right != nil {
    nn := pred(n)
    
    tmp := nn.Left
    nn.Left = n.Left
    n.Left = tmp
    if n.Left != nil {
      n.Left.Parent = n
    }
    if nn.Left != nil {
      nn.Left.Parent = nn
    }
    
    tmp = nn.Right
    nn.Right = n.Right
    n.Right = tmp
    if n.Right != nil {
      n.Right.Parent = n
    }
    if nn.Right != nil {
      nn.Right.Parent = nn
    }
    np := this.get_parent(n)
    nnp := this.get_parent(nn)
    tmp = nn.Parent
    nn.Parent = n.Parent
    n.Parent = tmp
    *np = nn
    *nnp = n
  }
  
  np := this.get_parent(n)
  if n.Left != nil {
    *np = n.Left
    n.Left.Parent = n.Parent
  } else {
    *np = n.Right
    if n.Right != nil {
      n.Right.Parent = n.Parent
    }
  }
  this.Size_--
}

func to_string(n* tree.Node)string {
  nn := dc(n)
  return fmt.Sprintf("n:%p p:%p c: (%v,%v)\n",nn,nn.Parent,nn.Key,nn.Value)
}

func (this* SPLAY) rotate_left(x *tree.Node)*tree.Node {
  if x.Right == nil {
    panic(fmt.Sprintf("rotate_left:x.Right==nil\n%s\n%d",this,this.Size_))
  }
  y := x.Right
  x.Right = y.Left
  if y.Left != nil {
    y.Left.Parent = x
  }
  y.Parent = x.Parent
  if x.Parent == nil {
    this.Root = y
  } else {
    if x == x.Parent.Left {
      x.Parent.Left = y
    } else {
      x.Parent.Right = y
    }
  }
  y.Left = x
  x.Parent = y
  return y
}

func (this* SPLAY) rotate_right(x *tree.Node) *tree.Node{
  if x.Left == nil {
    panic("rotate_right:x.Left==nil")
  }
  y := x.Left
  x.Left = y.Right
  if y.Right != nil {
    y.Right.Parent = x
  }
  y.Parent = x.Parent
  if x.Parent == nil {
    this.Root = y
  } else {
    if x == x.Parent.Left {
      x.Parent.Left = y
    } else {
      x.Parent.Right = y
    }
  }
  y.Right = x
  x.Parent = y
  return y;
}

func validate(node *Node)bool {
  if node == nil { return true }

  if node.Left != nil && !node.Left.Less(node.Key) {
    return false
  }
  if node.Right != nil && !node.Less(node.Right.Key) {
    return false
  }
  if node.Left != nil && node.Left.Parent != uc(node) {
    fmt.Println("node fail",node.Left.Key,node.Left.Parent.Key)
    return false
  }
  if node.Right != nil && node.Right.Parent != uc(node) {
    fmt.Println("node fail",node.Right.Key,node.Right.Parent.Key)
    return false
  }

  lh := validate(dc(node.Left))
  rh := validate(dc(node.Right))
  return lh && rh
}

func succ(n *tree.Node)*tree.Node {
  var nn *tree.Node
  for nn = n.Right ; nn.Left != nil ;{
    nn = nn.Left
  }
  return nn
}

func pred(n *tree.Node)*tree.Node {
  var nn *tree.Node
  for nn = n.Left ; nn.Right != nil ;{
    nn = nn.Right
  }
  return nn
}

func (this* SPLAY) get_parent(n* tree.Node)**tree.Node {
  if n.Parent == nil {
    return &this.Root
  }
  if n.Parent.Left == n {
    return &n.Parent.Left
  } else {
    return &n.Parent.Right
  }
}

func (this* SPLAY) Find(k tree.Key) (rc tree.Iterator) {
  var n,tmp = this.Root,(*tree.Node)(nil)
  var prev *tree.Node
  for n != nil {
    prev = n
    ret := k.Less(n.Key)
    if ret {
      n = n.Left
    } else {
      tmp = n
      n = n.Right
    }
  }
  if tmp != nil && !tmp.Less(k) {
    rc = tree.NewIter(tmp)
    this.splay(tmp)
    return
  }
  this.splay(prev)
  return
}
func (this* SPLAY) splay(n *tree.Node) {
  for n.Parent != nil {
    if n.Parent.Parent == nil {
      if n.Parent.Left == n {
        this.rotate_right(n.Parent)
      } else {
        this.rotate_left(n.Parent)
      }
    } else if n.Parent.Left == n && n.Parent.Parent.Left == n.Parent {
        this.rotate_right(n.Parent.Parent)
        this.rotate_right(n.Parent)
    } else if n.Parent.Right == n && n.Parent.Parent.Right == n.Parent {
        this.rotate_left(n.Parent.Parent)
        this.rotate_left(n.Parent)
    } else if n.Parent.Left == n && n.Parent.Parent.Right == n.Parent {
        this.rotate_right(n.Parent)
        this.rotate_left(n.Parent)
    } else {
        this.rotate_left(n.Parent)
        this.rotate_right(n.Parent)
    }
  }
}
