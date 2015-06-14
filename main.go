package main

import (
  "fmt"
  "github.com/bmaxa/trees/avl"
  "github.com/bmaxa/trees/treap"
  "github.com/bmaxa/trees/aa"
  "github.com/bmaxa/trees/rb"
  "github.com/bmaxa/trees/splay"
  "github.com/bmaxa/trees/tree"
  "math/rand"
  "time"
  . "testing"
  "runtime"
  "os"
  "runtime/debug"
  "strconv"
)

var N = 1000000
var  rnd1,rnd2,rnd3 []int

func init(){
  if len(os.Args)>1 {
    n,err := strconv.Atoi(os.Args[1])
    if err == nil {
      N = n
    }
  }
  fmt.Println(N)
  rnd1,rnd2,rnd3 = RandomShuffle(N),RandomShuffle(N),array(N)
}

//var  rnd1,rnd2 = array(N),RandomShuffle(N)

func RandomShuffle(n int)[]int{
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  fmt.Println("shuffle")
  return r.Perm(n)
}

func array(n int)[]int {
  arr := make([]int,n)
  for i:=0;i<n;i++{
    arr[i]=i
  }
  return arr
}

type Key int32
func (t Key) Less(k tree.Key) bool {
  return t<k.(Key)
}

func bench(t tree.ITree,rnd1,rnd2[]int,header string, prt bool) (inser,find,iter,ins,delet int64){
  fmt.Printf("%s\n--------------------------------\n",header)
  inser = Benchmark(func (b *B) {
  b.N = N
  t.Clear()
  b.ResetTimer()
  var item tree.Item
  for i,k := range rnd1 {
    item.Key = Key(k)
    item.Value = N-k
    t.Insert(item)
//    if prt { fmt.Println("inserting\n",t) }
    t.Delete(Key(rnd2[i]))
  }
  }).NsPerOp()
//  if prt { fmt.Println(t) }
  find = Benchmark(func (b *B) {
//  b.N = N
//  b.ResetTimer()
  for j,i:=0,0;j<b.N;j++ {
    if j%100 == 0 { i=0 }
    t.Find(Key(rnd1[i]))
    i++
  }
  }).NsPerOp()
  
  iter = Benchmark(func (b *B) {
  b.N = N
  b.ResetTimer()
  for i := t.Begin(); i != t.End(); i=i.Next() {
  }
  for i := t.RBegin(); i != t.End(); i=i.Prev() {
  }
  }).NsPerOp()
  left,right := t.Weight()
  fmt.Printf("height : %v\nweight : (%d,%d)\nvalid : %v\nsize %d\n",t.Height(),left,right,t.Validate(),t.Size())
  t.Clear()
  var printed bool 
  ins = Benchmark(func (b *B) {
  b.N = N
  t.Clear()
  b.ResetTimer()
  var item tree.Item
  for _,k := range rnd1 {
    item.Key = Key(k)
    item.Value = N-k
    t.Insert(item)
  }
  if !printed {
    fmt.Println("ascending iteration")
    for i := t.Find(Key(N-10)); i != t.End(); i=i.Next() {
      fmt.Println(i.Value())
    }
    fmt.Println("descending iteration")
    for i := t.Find(Key(10)); i != t.End(); i=i.Prev() {
      fmt.Println(i.Value())
    }
    printed = true
  }
  }).NsPerOp()
  left,right = t.Weight()
  fmt.Printf("height : %v\nweight : (%d,%d)\nvalid : %v\nsize %d\n",t.Height(),left,right,t.Validate(),t.Size())
  delet = Benchmark(func (b *B) {
  b.N = N
  var item tree.Item
  for _,k := range rnd1 {
    item.Key = Key(k)
    item.Value = N-k
    t.Insert(item)
  }
//  if prt { fmt.Println("deleting:\n",t) }
  b.ResetTimer()
  for _,k := range rnd1 {
    t.Delete(Key(k))
//    if prt { fmt.Println(t) }
  }
  }).NsPerOp()
  if prt { fmt.Println(t) }
  return
}

func main() {
  t := rb.New()
  debug.SetMaxStack(1<<32)
  rb_inser,rb_find,rb_iter,rb_ins,rb_delet := bench(t,rnd1,rnd2,"Red Black",false)
  tt := aa.New()
  aa_inser,aa_find,aa_iter,aa_ins,aa_delet := bench(tt,rnd1,rnd2,"AA",false)
  ttt := treap.New(nil)
  tr_inser,tr_find,tr_iter,tr_ins,tr_delet := bench(ttt,rnd1,rnd2,"Treap",false)
  tttt := avl.New()
  avl_inser,avl_find,avl_iter,avl_ins,avl_delet := bench(tttt,rnd1,rnd2,"AVL",true)
  ttttt := splay.New()
  splay_inser,splay_find,splay_iter,splay_ins,splay_delet := bench(ttttt,rnd1,rnd2,"SPLAY",true)
  ms := runtime.MemStats{}
  fmt.Printf("Tree\t\tInsert/Erase\tFind\t\tIter\t\tInsert\t\tDelete\n")
  fmt.Printf("SPLAY\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",splay_inser,splay_find,splay_iter,splay_ins,splay_delet)
  fmt.Printf("AVL\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",avl_inser,avl_find,avl_iter,avl_ins,avl_delet)
  fmt.Printf("RB\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",rb_inser,rb_find,rb_iter,rb_ins,rb_delet)
  fmt.Printf("AA\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",aa_inser,aa_find,aa_iter,aa_ins,aa_delet)
  fmt.Printf("TR\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",tr_inser,tr_find,tr_iter,tr_ins,tr_delet)
  rb_inser,rb_find,rb_iter,rb_ins,rb_delet = bench(t,rnd3,rnd2,"Red Black",false)
  aa_inser,aa_find,aa_iter,aa_ins,aa_delet = bench(tt,rnd3,rnd2,"AA",false)
  tr_inser,tr_find,tr_iter,tr_ins,tr_delet = bench(ttt,rnd3,rnd2,"Treap",false)
  avl_inser,avl_find,avl_iter,avl_ins,avl_delet = bench(tttt,rnd3,rnd2,"AVL",true)
  splay_inser,splay_find,splay_iter,splay_ins,splay_delet = bench(ttttt,rnd3,rnd2,"SPLAY",true)
  fmt.Println("Ordered Insert")
  fmt.Printf("Tree\t\tInsert/Erase\tFind\t\tIter\t\tInsert\t\tDelete\n")
  fmt.Printf("SPLAY\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",splay_inser,splay_find,splay_iter,splay_ins,splay_delet)
  fmt.Printf("AVL\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",avl_inser,avl_find,avl_iter,avl_ins,avl_delet)
  fmt.Printf("RB\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",rb_inser,rb_find,rb_iter,rb_ins,rb_delet)
  fmt.Printf("AA\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",aa_inser,aa_find,aa_iter,aa_ins,aa_delet)
  fmt.Printf("TR\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\n",tr_inser,tr_find,tr_iter,tr_ins,tr_delet)
  runtime.ReadMemStats(&ms)
  fmt.Printf("mallocs %v , frees %v\n%v\n",ms.Mallocs,ms.Frees,ms.BySize)
}
