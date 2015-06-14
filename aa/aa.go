package aa

import (
  "fmt"
  "github.com/bmaxa/trees/tree"
  "unsafe"
)

type AA struct {
  tree.Tree
}

type Node struct {
  tree.Node
  level uint
}

func uc(n *Node) *tree.Node{
  rc := (*tree.Node)(unsafe.Pointer(n))
  return rc
}
     
func dc(n *tree.Node)*Node{
  rc := (*Node)(unsafe.Pointer(n))
  return rc
}
          
func New()*AA {
  return &AA{}
}

func NewNode(i tree.Item)*tree.Node {
  n := tree.Node{ Item:i }
  return uc(&Node{ Node:n,level:1 })
}

func (t* AA) Validate()bool {
  return validate(dc(t.Root))
}

func (t AA) String() string {
  if t.Root == nil { return "Empty AA" }
  return tree.ToString(t.Root,"",true,to_string)
}
    
func to_string(n* tree.Node)string {
  nn := dc(n)
  return fmt.Sprintf("n:%p p:%p l:%d (%v,%v)\n",nn,nn.Parent,nn.level,nn.Key,nn.Value)
}

func skew(node **tree.Node) {
  if *node == nil || (*node).Left == nil || 
     dc(*node).level != dc((*node).Left).level {
    return
  }
  left := (*node).Left
  (*node).Left = left.Right
  if left.Right != nil { left.Right.Parent = *node }
  left.Right = *node
  left.Parent = (*node).Parent
  (*node).Parent = left
  *node = left
}

func split(node **tree.Node) {
  if *node == nil || (*node).Right == nil || (*node).Right.Right == nil ||
     dc(*node).level != dc((*node).Right.Right).level {
     return
  }
  right := (*node).Right
  (*node).Right = right.Left
  if right.Left != nil { right.Left.Parent = *node }
  right.Left = *node
  right.Parent = (*node).Parent
  (*node).Parent = right
  *node = right
  dc(*node).level++
}

func (this* AA) Insert(ins_value tree.Item) (tree.Value,bool) {
  if this.Root == nil {
    this.Root = NewNode(ins_value)
    this.Size_++
    return this.Root.Value,true
  }
  var node,prev,rc *tree.Node = this.Root,nil,nil
  var ret bool
  for node!=nil {
    prev = node
    ret = ins_value.Less(node.Key)
    if ret {
      node = node.Left
    } else {
      rc = node
      node = node.Right
    }
  }
  if rc != nil && !rc.Less(ins_value.Key) {
    return rc.Value,false
  }
  rc = NewNode(ins_value)
  if ret {
    prev.Left = rc
  } else {
    prev.Right = rc
  }
  rc.Parent = prev
  prev = rc
  for ;prev != nil;prev=prev.Parent{
    parent_node := this.get_parent(prev) 
    
    skew(parent_node)
    split(parent_node)
  }
  this.Size_++
  return rc.Value,true
}

func (this* AA) Delete(k tree.Key) bool {
  var node,prev,rc *tree.Node = this.Root,nil,nil
  var ret bool
  for node!=nil {
    prev = node
    ret = k.Less(node.Key)
    if ret {
      node = node.Left
    } else {
      rc = node
      node = node.Right
    }
  }
  if rc == nil || rc.Less(k) {
    return false
  }

  parent_node := this.get_parent(prev)

  rc.Item = prev.Item
  if prev.Right != nil {
    prev.Right.Parent = prev.Parent
    *parent_node = prev.Right
  } else {
    if prev.Left != nil {
      prev.Left.Parent = prev.Parent
    }
    *parent_node = prev.Left
  }
  
  for prev=prev.Parent;prev != nil;prev=prev.Parent {
    parent_node = this.get_parent(prev)
    pn := *parent_node
    var llevel,rlevel uint
    if pn.Left != nil {
      llevel = dc(pn.Left).level
    }
    if pn.Right != nil {
      rlevel = dc(pn.Right).level
    }
    should_be := min(llevel,rlevel)+1
    if should_be < dc(pn).level {
      dc(pn).level = should_be
      if pn.Right != nil && should_be < rlevel {
        dc(pn.Right).level = should_be
      }
    }
    skew(parent_node)
    skew(&(*parent_node).Right)
    if (*parent_node).Right != nil {
      skew(&(*parent_node).Right.Right)
    }
    split(parent_node)
    split(&(*parent_node).Right)
  }
  this.Size_--
  return true
}

func (this* AA) get_parent(prev *tree.Node)(parent_node **tree.Node) {
    if prev.Parent != nil {
      if prev.Parent.Left == prev {
        parent_node = &prev.Parent.Left
        } else {
        parent_node = &prev.Parent.Right
        }
    } else {
      parent_node = &this.Root
    }
    return
}

func min(a,b uint)uint {
  if a<b { 
    return a 
  } else {
    return b 
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
    fmt.Println("node fail",node.Left.Key,node.Left.Parent.Key)
    return false
  }
  if node.Right != nil && node.Right.Parent != uc(node) {
    fmt.Println("node fail",node.Right.Key,node.Right.Parent.Key)
    return false
  }
  if node.Left == nil && node.Right == nil && node.level != 1 {
    fmt.Println("leaf not 1",node.Key,node.level)
    return false
  }
  if node.Left != nil && node.level != dc(node.Left).level + 1 {
    fmt.Println("left node not less by 1",node.Key,node.level,node.Left.Key,dc(node.Left).level)
    return false
  }
  if node.Right != nil && 
     node.level != dc(node.Right).level && node.level != dc(node.Right).level + 1 {
    fmt.Println("right node not <= ",node.level,dc(node.Right).level)
    return false
  }
  if node.Left != nil && node.Right != nil && node.level == 1 {
    fmt.Println("node has two children and level  == 1",node.Key)
    return false
  }
  rc1 := validate(dc(node.Left))
  rc2 := validate(dc(node.Right))
  return rc1 && rc2
}
