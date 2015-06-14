package rb

import (
  "fmt"
  "github.com/bmaxa/trees/tree"
  "unsafe"
)

type RB struct{
  tree.Tree
}

const (
  RED = true
  BLACK = false
)

type Node struct {
  tree.Node
  colour bool
}

func uc(n *Node) *tree.Node{
  rc := (*tree.Node)(unsafe.Pointer(n))
  return rc
}
     
func dc(n *tree.Node)*Node{
  rc := (*Node)(unsafe.Pointer(n))
  return rc
}
          
func New()*RB{
  return &RB{}
}

func NewNode(i tree.Item,colour bool)*tree.Node {
  n := tree.Node{ Item:i }
  return uc(&Node{ Node:n , colour:colour })
}

func (t* RB) Validate()bool {
  return validate(dc(t.Root)) > 0
}

func (t RB) String() string {
  if t.Root == nil { return "Empty RB" }
  return tree.ToString(t.Root,"",true,to_string)
}

func (this* RB) Insert(i tree.Item)(tree.Value,bool) {
  x,rc := this.insert_helper(i)
  if !rc { return x.Value,rc }
  xx := x
  
  dc(x).colour = RED
  for x != this.Root && dc(x.Parent).colour == RED {
    if x.Parent == x.Parent.Parent.Left {
      y := x.Parent.Parent.Right
      if y != nil && dc(y).colour == RED {
        dc(x.Parent).colour = BLACK
        dc(y).colour = BLACK
        dc(x.Parent.Parent).colour = RED
        x = x.Parent.Parent
      } else {
        if x == x.Parent.Right {
          x = x.Parent
          this.rotate_left(x)
        }
        dc(x.Parent).colour = BLACK
        dc(x.Parent.Parent).colour = RED
        this.rotate_right(x.Parent.Parent)
      }
    } else {
      y := x.Parent.Parent.Left
      if y != nil && dc(y).colour == RED {
        dc(x.Parent).colour = BLACK
        dc(y).colour = BLACK
        dc(x.Parent.Parent).colour = RED
        x = x.Parent.Parent
      } else {
        if x == x.Parent.Left {
          x = x.Parent
          this.rotate_right(x)
        }
        dc(x.Parent).colour = BLACK
        dc(x.Parent.Parent).colour = RED
        this.rotate_left(x.Parent.Parent)
      }
    }
  }
  dc(this.Root).colour = BLACK
  return xx.Value,true
}

func (this* RB) Delete(k tree.Key)bool {
  n := this.Find(k).Node()
  if n == nil {
    return false
  } 
  this.delete(n)
  return true
}

func (this* RB) delete(z *tree.Node) {
  var x,y *tree.Node
  if z.Left == nil || z.Right == nil {
    y = z
  } else {
    y = tree.NewIter(z).Next().Node()
  }
  if y.Left == nil {
    x = y.Right
  } else {
    x = y.Left
  }
  if x != nil { x.Parent = y.Parent }
  if y.Parent == nil {
    this.Root = x
  } else {
    if y == y.Parent.Left {
      y.Parent.Left = x
    } else {
      y.Parent.Right = x
    }
  }
  if y != z {
    z.Item = y.Item
  }
  if dc(y).colour == BLACK {
    this.delete_fixup(x,y.Parent)
  }
  this.Size_--
}

func to_string(n* tree.Node)string {
  nn := dc(n)
  var colour string
  if nn.colour {
    colour = "Red"
  } else {
    colour = "Black"
  }
  return fmt.Sprintf("n:%p p:%p c:%s (%v,%v)\n",nn,nn.Parent,colour,nn.Key,nn.Value)
}

func (this* RB) rotate_left(x *tree.Node) {
  if x.Right == nil {
    panic("rotate_left:x.Right==nil")
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
}

func (this* RB) rotate_right(x *tree.Node) {
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
}

func (this* RB) insert_helper(i tree.Item)(*tree.Node,bool) {
  var x,y,z *tree.Node = this.Root,nil,nil
  if x == nil { 
    z = NewNode(i,RED)
    this.Root = z
    this.Size_++
    return z,true
  }
  
  for x != nil {
    y = x
    if i.Less(x.Key) {
      x = x.Left
    } else {
      z = x
      x = x.Right
    }
  }
  if z != nil && !z.Less(i.Key) {
    return z,false
  }
  z = NewNode(i,RED)
  z.Parent = y
  if z.Less(y.Key) {
    y.Left = z
  } else {
    y.Right = z
  }
  this.Size_++
  return z,true
}

func (this* RB) delete_fixup(x,p *tree.Node) {
  for x != this.Root && !is_red(dc(x)) {
    if x == p.Left {
      w := p.Right
      if is_red(dc(w)) {
        dc(w).colour = BLACK
        dc(p).colour = RED
        this.rotate_left(p)
        w = p.Right
      }
      if !is_red(dc(w.Left)) && !is_red(dc(w.Right)) {
        dc(w).colour = RED
        x = p
        p = p.Parent
      } else {
        if !is_red(dc(w.Right)) {
          dc(w.Left).colour = BLACK
          dc(w).colour = RED
          this.rotate_right(w)
          w = p.Right
        }
        dc(w).colour = dc(p).colour
        dc(p).colour = BLACK
        dc(w.Right).colour = BLACK
        this.rotate_left(p)
        p = this.Root
        x = p
      }
    } else {
      w := p.Left
      if is_red(dc(w)) {
        dc(w).colour = BLACK
        dc(p).colour = RED
        this.rotate_right(p)
        w = p.Left
      }
      if !is_red(dc(w.Right)) && !is_red(dc(w.Left)) {
        dc(w).colour = RED
        x = p
        p = p.Parent
      } else {
        if !is_red(dc(w.Left)) {
          dc(w.Right).colour = BLACK
          dc(w).colour = RED
          this.rotate_left(w)
          w = p.Left
        }
        dc(w).colour = dc(p).colour
        dc(p).colour = BLACK
        dc(w.Left).colour = BLACK
        this.rotate_right(p)
        p = this.Root
        x = p
      }
    }
  }
  if x != nil { dc(x).colour = BLACK }
}

func validate(node *Node)int {
  if node == nil { return 1 }
  if is_red(node) {
    if is_red(dc(node.Left)) || is_red(dc(node.Right)) {
      fmt.Println("red violation")
      return 0
    }
  }
  lh := validate(dc(node.Left))
  rh := validate(dc(node.Right))

  if node.Left != nil && !node.Left.Less(node.Key) {
    return 0
  }
  if node.Right != nil && !node.Less(node.Right.Key) {
    return 0
  }
  if node.Left != nil && node.Left.Parent != uc(node) {
    fmt.Println("node fail",node.Left.Key,node.Left.Parent.Key)
    return 0
  }
  if node.Right != nil && node.Right.Parent != uc(node) {
    fmt.Println("node fail",node.Right.Key,node.Right.Parent.Key)
    return 0
  }
  if lh != 0 && rh != 0 && lh != rh {
    panic(fmt.Sprintf ( "Black violation :\n%s\n",tree.ToString(uc(node),"",true,to_string) ))
    return 0
  }  
  if lh != 0 && rh != 0 {
    if is_red(node) {
      return lh
    } else {
      return lh+1
    }
  }
  return 0
}

func is_red(node *Node) bool {
  return node != nil && node.colour == RED
}
