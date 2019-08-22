package relgom

// // Iterator returns an iterator over Pet tuples in r.
// func (r PetRelation) Iterator() PetIterator {
// 	return &petIterator{r.model, r.set, nil}
// }

// // PetIterator provides for iteration over a set of Pet tuples.
// type PetIterator interface {
// 	MoveNext() bool
// 	Current() *Pet
// }

// type petIterator struct {
// 	model PetShopModel
// 	set   *seq.HashMap
// 	t     *Pet
// }

// // MoveNext implements seq.Setable.
// func (i *petIterator) MoveNext() bool {
// 	kv, set, has := i.set.FirstRestKV()
// 	if has {
// 		i.set = set
// 		i.t = &Pet{petData: kv.Val.(*petData), model: i.model}
// 	}
// 	return has
// }

// // Current implements seq.Setable.
// func (i *petIterator) Current() *Pet {
// 	if i.t == nil {
// 		panic("no current Pet")
// 	}
// 	return i.t
// }
