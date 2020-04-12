package graph

import (
	"sort"

	"github.com/Tom-Johnston/mamba/ints"
	"github.com/Tom-Johnston/mamba/sortints"
)

//IsPlanar returns true if the graph is a planar graph and false otherwise. The current implementation is simple but runs in O(n^2) while linear time algorithms are known.
//TODO Which algorithm does this use?
//TODO Can we make it output a planar embedding?
func IsPlanar(g Graph) bool {
	if g.N() < 5 {
		return true
	}

	//Find the biconnected components and then perform the count on each one.
	//fmt.Println(g)
	bicoms, _ := BiconnectedComponents(g)
	//fmt.Println(bicoms)
	for _, bicom := range bicoms {
		h := InducedSubgraph(g, bicom)
		//fmt.Println(h)
		n := h.N()
		if n < 5 {
			continue
		}
		if h.M() > 3*n-6 {
			return false
		}

		//TODO trees
		//Find any cycle.
		//TODO The algorithm works with any cycle but do we want to choose one anyway?
		parents := make([]int, n)
		for i := 1; i < n; i++ {
			parents[i] = -1
		}
		////fmt.Println(parents)
		//TODO: Keep track of depth as in fundCycles
		toCheck := make([]int, 1)
		var v int
		HV := make([]int, 0, n)
		HE := make([]int, 0, h.M())
	findCycle:
		for len(toCheck) > 0 {
			v = toCheck[len(toCheck)-1]
			//fmt.Println(v)
			for _, u := range h.Neighbours(v) {
				//fmt.Println(neighbourhoods[v])
				if u == parents[v] {
					continue
				}
				// //fmt.Println("u", u)
				//fmt.Println(parents)
				if parents[u] != -1 {
					//fmt.Println("u", u)
					//fmt.Println("v", v)
					HV = append(HV, u)

					HV = append(HV, v)
					if u < v {
						HE = append(HE, (v*(v-1))/2+u)
					} else {
						HE = append(HE, (u*(u-1))/2+v)
					}
					//fmt.Println("H1E", HE)
					previous := v
					for {
						var e int
						//fmt.Println(previous, parents[previous])
						if previous < parents[previous] {
							e = (parents[previous]*(parents[previous]-1))/2 + previous
						} else {
							e = (previous*(previous-1))/2 + parents[previous]
						}
						previous = parents[previous]

						// //fmt.Println(parents[u])
						HE = append(HE, e)
						if previous == u {
							break findCycle
						}

						HV = append(HV, previous)
						//fmt.Println("HE2", HE)
					}
				}
				toCheck = append(toCheck, u)
				parents[u] = v
				continue findCycle
			}
			toCheck = toCheck[:len(toCheck)-1]
		}

		//fmt.Println("Cycle found")
		HF := make([][]int, 2) //The set of faces of H
		HF[0] = make([]int, len(HV))
		copy(HF[0], HV)
		HF[1] = make([]int, len(HV))
		copy(HF[1], HV)

		// //fmt.Println("HF", HF)
		// //fmt.Println(HV)
		// //fmt.Println(HE)
		sort.Ints(HV)
		sort.Ints(HE)
		//fmt.Println(HE)

		//Find the H-fragments
		type frag struct {
			E []int //Edges in the H fragment
			V []int //Vertices in the H fragment excluding attachments
			F []int //The list of faces that H can currently be embedded in.
			A []int //The attachment vertices
		}

		HFrags := make([]frag, 0)

		//Find the internal H-fragments
		for _, v = range HV {
			for _, u := range h.Neighbours(v) {
				if u > v {
					continue
				}
				index := sort.SearchInts(HV, u)
				if index < len(HV) && HV[index] == u {
					var e int
					if u < v {
						e = (v*(v-1))/2 + u
					} else {
						e = (u*(u-1))/2 + v
					}
					indexE := sort.SearchInts(HE, e)
					if indexE >= len(HE) || HE[indexE] != e {
						//This edge forms a H fragment.
						if u < v {
							f := frag{E: []int{e}, V: nil, F: []int{0, 1}, A: []int{u, v}}
							HFrags = append(HFrags, f)
						} else {
							f := frag{E: []int{e}, V: nil, F: []int{0, 1}, A: []int{v, u}}
							HFrags = append(HFrags, f)
						}
					}
				}
			}
		}

		//fmt.Println("Hfrags", HFrags)

		//Find the external H-fragments.

		unseen := sortints.Complement(n, HV)
		////fmt.Println("unseen", unseen)
		seen := make([]int, 0, h.N())
		attach := make([]int, 0, 2)
		E := make([]int, 0, 2)
		for len(unseen) > 0 {
			seen = seen[:1]
			toCheck = toCheck[:1]
			toCheck[0], unseen = unseen[len(unseen)-1], unseen[:len(unseen)-1]
			seen[0] = toCheck[0]
			attach = attach[:0]
			E = E[:0]
			for len(toCheck) > 0 {
				v, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]
				for _, u := range h.Neighbours(v) {
					indexU := sort.SearchInts(unseen, u)
					////fmt.Println("u", u, "indexU", indexU)
					if indexU < len(unseen) && unseen[indexU] == u {
						////fmt.Println("unseen")
						//We haven't seen this vertex before.
						//Remove the vertex from unseen.
						unseen = unseen[:indexU+copy(unseen[indexU:], unseen[indexU+1:])]
						//Add the vertex to seen.
						//TODO Insert this properly.
						seen = append(seen, u)
						sort.Ints(seen)
						//Add the edge to E;
						if u < v {
							e := (v*(v-1))/2 + u
							E = append(E, e)
						} else {
							e := (u*(u-1))/2 + v
							E = append(E, e)
						}
						//Add the vertex to toCheck
						toCheck = append(toCheck, u)
						continue
					}

					indexS := sort.SearchInts(seen, u)
					if indexS < len(seen) && seen[indexS] == u {
						////fmt.Println("seen")
						//We have seen this vertex before. No need to remove from unseen or add to seen. Just add the edge to E.
						//Add the edge to E;
						if u < v {
							e := (v*(v-1))/2 + u
							E = append(E, e)
						} else {
							e := (u*(u-1))/2 + v
							E = append(E, e)
						}
						continue
					}

					//The vertex must be an attachement point. Add the vertex to the list of attachments and add the edge to E.
					//Add the vertex to the list of attachments.
					indexA := sort.SearchInts(attach, u)
					if indexA >= len(attach) || attach[indexA] != u {
						attach = append(attach, u)
						sort.Ints(attach)
					}

					//Add the edge to E;
					if u < v {
						e := (v*(v-1))/2 + u
						E = append(E, e)
					} else {
						e := (u*(u-1))/2 + v
						E = append(E, e)
					}
				}
			}
			//We have finished exploring this HFragment I think. Week long break not helpful...

			tmpV := make([]int, len(seen))
			copy(tmpV, seen)
			sort.Ints(tmpV)
			tmpE := make([]int, len(E))
			copy(tmpE, E)
			sort.Ints(tmpE)
			tmpA := make([]int, len(attach))
			copy(tmpA, attach)
			sort.Ints(tmpA)
			HFrags = append(HFrags, frag{E: tmpE, V: tmpV, F: []int{0, 1}, A: tmpA})
		}
		////fmt.Println(HFrags)
		//TODO Use a heap to maintain this instead of looping over all fragments.
	fragLoop:
		for len(HFrags) > 0 {
			//fmt.Println(HFrags)
			//fmt.Println(HF)
			var f frag
			oldSize := len(HFrags)
			for i, f2 := range HFrags {
				if len(f2.F) == 1 {
					f = f2
					HFrags[i] = HFrags[len(HFrags)-1]
					HFrags = HFrags[:len(HFrags)-1]
					break
				}
			}
			if len(f.E) == 0 {
				//We haven't selected a fragment.
				f, HFrags = HFrags[len(HFrags)-1], HFrags[:len(HFrags)-1]
			}
			//fmt.Println("f", f)
			//Find an alpha path by DFS
			toCheck = toCheck[:1]
			seen = seen[:1]
			toCheck[0] = f.A[0]
			seen[0] = toCheck[0]
			//TODO Don't need all of these...
			parents := make([]int, n)
			for i := range parents {
				parents[i] = -1
			}
		outerLoop:
			for { //We must find a path so len(toCheck) > 0 unnecessary
				//fmt.Println(toCheck)
				////fmt.Println(f)
				if len(toCheck) == 0 {
					//fmt.Println(f)
					panic("Oh dear")
				}
				v = toCheck[len(toCheck)-1]
				for j := 0; j < len(f.V); j++ {
					//TODO Replace this with unseen?
					////fmt.Println(j)
					//fmt.Println(v)
					//fmt.Println("f.V[j]", f.V[j])
					if parents[f.V[j]] == -1 && isEdge(f.V[j], v, f.E) {
						////fmt.Println("add")
						toCheck = append(toCheck, f.V[j])
						//fmt.Println(toCheck)
						parents[f.V[j]] = v
						continue outerLoop
					}
				}
				//fmt.Println("j Loop")
				//fmt.Println("parents", parents)
				for j := 1; j < len(f.A); j++ {
					if isEdge(f.A[j], v, f.E) {
						//We have a path
						//Extract the path.
						attachA := f.A[j]
						attachB := 0
						//fmt.Println(attachB)
						seenVertices := make([]int, 0, 2)
						parents[f.A[j]] = v
						parent := v
						v = f.A[j]
						for {
							//fmt.Println("parent", parent, "v", v)
							seenVertices = append(seenVertices, v)
							if parent < v {
								e := (v*(v-1))/2 + parent
								HE = append(HE, e)
							} else {
								e := (parent*(parent-1))/2 + v
								HE = append(HE, e)
							}
							if parents[parent] == -1 {
								attachB = parent
								seenVertices = append(seenVertices, parent)
								break
							}
							v = parent
							parent = parents[v]
						}

						HV = append(HV, seenVertices[1:]...)
						sort.Ints(HV)
						//fmt.Println("HV", HV)
						sort.Ints(HE)
						// //fmt.Println("HE", HE)
						// //fmt.Println("HV", HV)
						// //fmt.Println("HF", HF)
						//Update the faces.

						faceIndex := f.F[len(f.F)-1]
						face := HF[faceIndex]
						////fmt.Println(face)
						newFaceA := make([]int, 0, len(face))
						newFaceB := make([]int, 0, len(face))
						inA := true
						//fmt.Println("Attach", attachA, attachB)
						//fmt.Println("face", face)
						//fmt.Println(seenVertices)
						for _, w := range face {
							//fmt.Println("w", w)
							//fmt.Println(newFaceA, newFaceB)
							if inA {
								newFaceA = append(newFaceA, w)
							} else {
								newFaceB = append(newFaceB, w)
							}
							if w == attachA || w == attachB {
								inA = !inA
								if inA {
									newFaceA = append(newFaceA, w)
								} else {
									newFaceB = append(newFaceB, w)
									if w == attachA {
										newFaceA = append(newFaceA, seenVertices[1:len(seenVertices)-1]...)
									} else {
										newFaceA = append(newFaceA, ints.Reverse(seenVertices[1:len(seenVertices)-1])...)
										ints.Reverse(seenVertices[1 : len(seenVertices)-1])
									}
								}
							}
						}

						//Add the path to face B. Got to make sure the orientation of the path is correct.
						if newFaceB[0] == attachB {
							newFaceB = append(newFaceB, seenVertices[1:len(seenVertices)-1]...)
						} else {
							newFaceB = append(newFaceB, ints.Reverse(seenVertices[1:len(seenVertices)-1])...)
						}
						//fmt.Println(newFaceA, newFaceB)
						//Add the new faces to the list of HF.
						HF[faceIndex] = newFaceA
						HF = append(HF, newFaceB)

						// //fmt.Println("HF", HF)

						//TODO Find the new HFragments

						//Find the new internal H-fragments
						for _, v = range seenVertices[1 : len(seenVertices)-1] {
							for _, u := range h.Neighbours(v) {
								index := sort.SearchInts(HV, u)
								if index < len(HV) && HV[index] == u {
									var e int
									if u < v {
										e = (v*(v-1))/2 + u
									} else {
										e = (u*(u-1))/2 + v
									}
									indexE := sort.SearchInts(HE, e)
									if indexE >= len(HE) || HE[indexE] != e {
										//This edge forms a H fragment.
										//Find the faces that it can be put in.
										//TODO Use a buffer and copy here?
										tmpF := make([]int, 0, 1)
										for j, face := range HF {
											uSeen := false
											vSeen := false
											for _, a := range face {
												if a == u {
													uSeen = true
												} else if a == v {
													vSeen = true
												}
												if uSeen && vSeen {
													tmpF = append(tmpF, j)
													break
												}
											}
										}
										if len(tmpF) == 0 {
											return false
										}
										if u < v {
											f := frag{E: []int{e}, V: nil, F: tmpF, A: []int{u, v}}
											HFrags = append(HFrags, f)
										} else {
											f := frag{E: []int{e}, V: nil, F: tmpF, A: []int{v, u}}
											HFrags = append(HFrags, f)
										}
									}
								}
							}
						}

						//fmt.Println("Hfrags", HFrags)
						//fmt.Println("Done")

						//Find the external H-fragments.

						unseen := f.V
						for _, v := range seenVertices[1 : len(seenVertices)-1] {
							indexV := sort.SearchInts(unseen, v)
							if indexV < len(unseen) && unseen[indexV] == v {
								unseen = unseen[:indexV+copy(unseen[indexV:], unseen[indexV+1:])]
							} else {
								panic("This shouldn't happen...")
							}
						}
						////fmt.Println("unseen", unseen)
						seen := make([]int, 0, h.N())
						attach := make([]int, 0, 2)
						E := make([]int, 0, 2)
						for len(unseen) > 0 {
							seen = seen[:1]
							toCheck = toCheck[:1]
							toCheck[0], unseen = unseen[len(unseen)-1], unseen[:len(unseen)-1]
							seen[0] = toCheck[0]
							attach = attach[:0]
							E = E[:0]
							for len(toCheck) > 0 {
								v, toCheck = toCheck[len(toCheck)-1], toCheck[:len(toCheck)-1]
								for _, u := range h.Neighbours(v) {
									////fmt.Println(seen)
									//fmt.Println("V", v, "u", u)
									indexU := sort.SearchInts(unseen, u)
									////fmt.Println("u", u, "indexU", indexU)
									if indexU < len(unseen) && unseen[indexU] == u {
										////fmt.Println("unseen")
										//We haven't seen this vertex before.
										//Remove the vertex from unseen.
										unseen = unseen[:indexU+copy(unseen[indexU:], unseen[indexU+1:])]
										//Add the vertex to seen.
										//TODO Insert this properly.
										seen = append(seen, u)
										sort.Ints(seen)
										//Add the edge to E;
										if u < v {
											e := (v*(v-1))/2 + u
											E = append(E, e)
										} else {
											e := (u*(u-1))/2 + v
											E = append(E, e)
										}
										//Add the vertex to toCheck
										toCheck = append(toCheck, u)
										continue
									}

									indexS := sort.SearchInts(seen, u)
									if indexS < len(seen) && seen[indexS] == u {
										////fmt.Println("seen")
										//We have seen this vertex before. No need to remove from unseen or add to seen. Just add the edge to E.
										//Add the edge to E;
										if u < v {
											e := (v*(v-1))/2 + u
											E = append(E, e)
										} else {
											e := (u*(u-1))/2 + v
											E = append(E, e)
										}
										continue
									}

									//The vertex must be an attachement point. Add the vertex to the list of attachments and add the edge to E.
									//Add the vertex to the list of attachments.
									indexA := sort.SearchInts(attach, u)
									if indexA >= len(attach) || attach[indexA] != u {
										attach = append(attach, u)
										sort.Ints(attach)
									}
									//Add the edge to E;
									if u < v {
										e := (v*(v-1))/2 + u
										E = append(E, e)
									} else {
										e := (u*(u-1))/2 + v
										E = append(E, e)
									}
								}
							}
							//We have finished exploring this HFragment I think. Week long break not helpful...

							tmpV := make([]int, len(seen))
							copy(tmpV, seen)
							sort.Ints(tmpV)
							tmpE := make([]int, len(E))
							copy(tmpE, E)
							sort.Ints(tmpE)
							tmpA := make([]int, len(attach))
							sort.Ints(tmpA)
							copy(tmpA, attach)
							tmpF := make([]int, 0)
							for j, face := range HF {
								if contains(face, tmpA) {
									tmpF = append(tmpF, j)
								}
							}
							if len(tmpF) == 0 {
								return false
							}
							HFrags = append(HFrags, frag{E: tmpE, V: tmpV, F: tmpF, A: tmpA})
						}

						//Update the faces of the other HFragments.
						for i := 0; i < oldSize-1; i++ {
							//fmt.Println(HFrags[i], faceIndex)
							indexF := sort.SearchInts(HFrags[i].F, faceIndex)
							if indexF < len(HFrags[i].F) && HFrags[i].F[indexF] == faceIndex {
								//fmt.Println("found")
								//The fragment could previously be embedded in this face so this needs updating.
								if !contains(newFaceA, HFrags[i].A) {
									HFrags[i].F[indexF] = HFrags[i].F[len(HFrags[i].F)-1]
									HFrags[i].F = HFrags[i].F[:len(HFrags[i].F)-1]
									sort.Ints(HFrags[i].F)
								}
								if contains(newFaceB, HFrags[i].A) {
									HFrags[i].F = append(HFrags[i].F, len(HF)-1)
								}
							}
							if len(HFrags[i].F) == 0 {
								return false
							}
						}
						//fmt.Println("faces", HF)
						//fmt.Println("Hfrags", HFrags)
						//fmt.Println("Done")
						continue fragLoop
					}
				}
				toCheck = toCheck[:len(toCheck)-1]
			}
		}
	}

	return true
}

func isEdge(u, v int, E []int) bool {
	if u == v {
		return false
	}
	var e int
	if u < v {
		e = (v*(v-1))/2 + u
	} else {
		e = (u*(u-1))/2 + v
	}
	indexE := sort.SearchInts(E, e)
	if indexE < len(E) && E[indexE] == e {
		return true
	}
	return false
}

//Checks if the unsorted slice a contains the sorted array b. There can't be any duplicates.
func contains(a, b []int) bool {
	numberNotSeen := len(b)
	for _, v := range a {
		indexB := sort.SearchInts(b, v)
		if indexB < len(b) && b[indexB] == v {
			numberNotSeen--
			if numberNotSeen == 0 {
				return true
			}
		}
	}
	return false
}
