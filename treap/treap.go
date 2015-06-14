package treap

import (
  "fmt"
  "math/rand"
  "time"
  "github.com/bmaxa/trees/tree"
  "unsafe"
)

type Treap struct {
  tree.Tree
  r *rand.Rand
}

type Node struct {
  tree.Node
  priority float64
}

func uc(n *Node) *tree.Node{
  rc := (*tree.Node)(unsafe.Pointer(n))
  return rc
}

func dc(n *tree.Node)*Node{
  rc := (*Node)(unsafe.Pointer(n))
  return rc
}

func New(r *rand.Rand) *Treap {
  rc := Treap{r:r}
  if r == nil {
    rc.r = rand.New(rand.NewSource(time.Now().UnixNano()))
  }
  return &rc
}

func (this* Treap) prn()float64 {
  return this.r.Float64()
}

func NewNode(i tree.Item,t* Treap) *tree.Node {
  n := tree.Node{Item:i}
  rc := uc(&Node{Node:n,priority:t.prn()})
  return rc
}

func (this* Treap) Insert(i tree.Item)(tree.Value,bool) {
  if this.Root == nil {
    this.Root = NewNode(i,this)
    this.Size_++
    return this.Root.Value,true
  }
  n := this.Root
  var rc,prev *tree.Node
  var ret bool
  for n != nil {
    prev = n
    ret = i.Less(n.Key)
    if ret {
      n = n.Left
    } else {
      rc = n
      n = n.Right
    }
  }
  
  if rc != nil && !rc.Less(i.Key) {
    return rc.Value,false
  }

  if ret {
    prev.Left = NewNode(i,this)
    rc = prev.Left
    prev.Left.Parent = prev
  } else {
    prev.Right = NewNode(i,this)
    rc = prev.Right
    prev.Right.Parent = prev
  }
  n = rc
  this.rebalance_up(dc(n))
  this.Size_++
  
  return rc.Value,true
}

func (this* Treap) Delete(k tree.Key)bool{
  rc,n := (*tree.Node)(nil),this.Root
  for n != nil {
    ret := k.Less(n.Key)
    if ret {
      n = n.Left
    } else {
      rc = n
      n = n.Right
    }
  }
  if rc == nil || rc.Less(k) {
    return false
  }
  
  var reb *Node
  for rc.Left != nil && rc.Right != nil {
    n := dc(rc).rotate_right()
    if this.Root == rc {
      this.Root = uc(n)
    }
    if reb == nil && n.Left != nil && dc(n.Left).priority<n.priority {
      reb = n
    }
  }
  var parent_node **tree.Node
  if rc.Parent != nil {
    if rc.Parent.Left == rc {
      parent_node = &rc.Parent.Left
    } else {
      parent_node = &rc.Parent.Right
    }
  }
  if rc.Left == nil && rc.Right == nil {
    if parent_node != nil {
      *parent_node = nil
    } else {
      this.Root = nil
    }
  }	else if rc.Left == nil {
    if parent_node != nil {
      *parent_node = rc.Right
      rc.Right.Parent = rc.Parent
    } else {
      this.Root = rc.Right
      rc.Right.Parent = nil
    }
    rc.Right = nil
  } else if rc.Right == nil {
    if parent_node != nil {
      *parent_node = rc.Left
      rc.Left.Parent = rc.Parent
    } else {
      this.Root = rc.Left
      rc.Left.Parent = nil
    }
    rc.Left = nil
  }
  this.Size_--
  this.rebalance_left(reb)
  return true
}

func (t* Treap) Validate()bool {
  return validate(dc(t.Root))
}

func (t Treap) String() string {
  if t.Root == nil { return "Empty Treap" }
  return tree.ToString(t.Root,"",true,to_string)
}

func to_string(n* tree.Node)string {
  nn := dc(n)
  return fmt.Sprintf("p:%f (%v,%v)\n",nn.priority,nn.Key,nn.Value)
}

func (this* Node) rotate_left()*Node {
  n := this.Left
  if n == nil { return nil }
  
  this.Left = n.Right
  if this.Left != nil {
    this.Left.Parent = uc(this)
  }
  
  n.Right = uc(this)
  
  if this.Parent != nil {
    if this.Parent.Left == uc(this) {
      this.Parent.Left = n
      n.Parent = this.Parent
    } else {
      if this.Parent.Right != uc(this) {
        s := fmt.Sprintf("rotate Left failed, child\n%s\nParent\n%s\n",this,this.Parent)
        panic(s)
      }
      this.Parent.Right = n
      n.Parent = this.Parent
    }
  } else {
    n.Parent = nil
  }
  this.Parent = n
  return dc(n)
}

func (this* Node) rotate_right()*Node {
  n := this.Right
  if n == nil { return nil }
  
  this.Right = n.Left
  if this.Right != nil {
    this.Right.Parent = uc(this)
  }
  
  n.Left = uc(this)
  
  if this.Parent != nil {
    if this.Parent.Left == uc(this) {
      this.Parent.Left = n
      n.Parent = this.Parent
    } else {
      if this.Parent.Right != uc(this) {
        s := fmt.Sprintf("rotate Right failed, child\n%s\nParent\n%s\n",this,this.Parent)
        panic(s)
      }
      this.Parent.Right = n
      n.Parent = this.Parent
    }
  } else {
    n.Parent = nil
  }
  this.Parent = n
  return dc(n)
}
func (this* Treap) rebalance_left(node *Node){
  if node == nil { return }
  this.rebalance_left(dc(node.Left))
  for node.Left != nil && dc(node.Left).priority<node.priority {
    n := node.rotate_left()
    if uc(node) == this.Root {
      this.Root = uc(n)
    }
  }
}

func (this* Treap) rebalance_right(node *Node){
  if node == nil { return }
  this.rebalance_right(dc(node.Right))
  for node.Right != nil && dc(node.Right).priority<node.priority {
    n := node.rotate_right()
    if uc(node) == this.Root {
      this.Root = uc(n)
    }
  }
}

func (this* Treap) rebalance_up(n* Node) {
  for n.Parent != nil && n.priority<dc(n.Parent).priority {
    if n.Parent.Left == uc(n) {
      n = dc(n.Parent).rotate_left()
    } else {
      n = dc(n.Parent).rotate_right()
    }
    if n.Parent == nil {
      this.Root = uc(n)
    }
  }
}

func validate(node* Node)bool {
  if node == nil { return true }
  if node.Left != nil && !node.Left.Less(node.Key) {
    return false
  }
  if node.Right != nil && !node.Less(node.Right.Key) {
    return false
  }
  if node.Left != nil && node.Left.Parent != uc(node) {
    return false
  }
  if node.Right != nil && node.Right.Parent != uc(node) {
    return false
  }
  if node.Left != nil && node.priority > dc(node.Left).priority {
    return false
  }
  if node.Right != nil && node.priority > dc(node.Right).priority {
    return false
  }
  rc1 := validate(dc(node.Left))
  rc2 := validate(dc(node.Right))
  return rc1 && rc2
}
