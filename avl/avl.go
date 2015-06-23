package avl

import (
  "fmt"
  "github.com/bmaxa/trees/tree"
  "unsafe"
)

type AVL struct{
  tree.Tree
}

type Node struct {
  tree.Node
  height int32
}

func uc(n *Node) *tree.Node{
  rc := (*tree.Node)(unsafe.Pointer(n))
  return rc
}
     
func dc(n *tree.Node)*Node{
  rc := (*Node)(unsafe.Pointer(n))
  return rc
}
          
func New()*AVL{
  return &AVL{}
}

func NewNode(i tree.Item)*tree.Node {
  n := tree.Node{ Item:i }
  return uc(&Node{ Node:n , height:1 })
}

func (t* AVL) Validate()bool {
  return validate(dc(t.Root))
}

func (t AVL) String() string {
  if t.Root == nil { return "Empty AVL" }
  return tree.ToString(t.Root,"",true,to_string)
}

func (this* AVL) Insert(i tree.Item)(tree.Value,bool) {
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
  
  for nn := prev;nn != nil ; nn = nn.Parent {
    this.balance(nn)
  }
  return n.Value,true
}

func (this* AVL) InsertEqual(i tree.Item)(tree.Value) {
  if this.Root == nil {
    this.Root = NewNode(i)
    this.Size_++
    return this.Root.Value
  }
  var prev *tree.Node
  var ret bool
  for n := this.Root ; n != nil; {
    prev = n
    if ret = i.Less(n.Key);ret {
      n = n.Left
    } else {
      n = n.Right
    }
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
  
  for nn := prev;nn != nil ; nn = nn.Parent {
    this.balance(nn)
  }
  return n.Value
}

func (this* AVL) Delete(k tree.Key)bool {
  n := this.Find(k).Node()
  if n == nil {
    return false
  } 
  this.delete(n)
  return true
}

func (this* AVL) DeleteIter(i tree.Iterator)bool {
  n := i.Node()
  if n == nil {
    return false
  }
  
  this.delete(n)
  return true
}

func (this* AVL) delete(n *tree.Node) {
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
  prev := n.Parent
  if n.Left != nil {
    *np = n.Left
    n.Left.Parent = n.Parent
  } else {
    *np = n.Right
    if n.Right != nil {
      n.Right.Parent = n.Parent
    }
  }
  for ;prev!=nil;prev=prev.Parent {
    this.balance(prev)
  }
  this.Size_--
}

func to_string(n* tree.Node)string {
  nn := dc(n)
  return fmt.Sprintf("n:%p p:%p h:%v (%v,%v)\n",nn,nn.Parent,nn.height,nn.Key,nn.Value)
}

func (this* AVL) rotate_left(x *tree.Node) {
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
  fixheight(x)
  fixheight(y)
}

func (this* AVL) rotate_right(x *tree.Node) {
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
  fixheight(x)
  fixheight(y)
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
  if abs(bfactor(uc(node))) > 1 {
    fmt.Println("abs(balance factor) > 1",node)
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

func (this* AVL) get_parent(n* tree.Node)**tree.Node {
  if n.Parent == nil {
    return &this.Root
  }
  if n.Parent.Left == n {
    return &n.Parent.Left
  } else {
    return &n.Parent.Right
  }
}

func height(n *tree.Node)int32 {
  nn := dc(n)
  if nn != nil {
    return nn.height
  } else {
    return 0
  }
}

func bfactor(n* tree.Node)int32 {
  return height(n.Right)-height(n.Left)
}

func fixheight(n* tree.Node) {
  l := height(n.Left)
  r := height(n.Right)
  if l > r {
    dc(n).height = l + 1
  } else {
    dc(n).height = r + 1
  }
}

func (this* AVL) balance(n* tree.Node) {
  fixheight(n)
  if bfactor(n) == 2 {
    if bfactor(n.Right) < 0 {
      this.rotate_right(n.Right)
    }
    this.rotate_left(n)
    return
  }
  if bfactor(n) == -2 {
    if bfactor(n.Left) > 0 {
      this.rotate_left(n.Left)
    }
    this.rotate_right(n)
    return
  }
}

func abs(v int32)int32 {
  if v < 0 { 
    return -v 
  } else {
    return v
  }
}
