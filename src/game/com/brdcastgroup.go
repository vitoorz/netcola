package com

import (
    . "types"
    "sync"
)

type BrdCastGroup struct {
    Id       IdString
    Members  map[IdString]int32 //index

    lock     sync.Mutex
}

func NewBrdCastGroup(creator IdString) *BrdCastGroup {
    return &BrdCastGroup{
        Id:      creator,
        Members: make(map[IdString]int32, 0),
    }
}

func (bcg *BrdCastGroup) GroupDetail() *BrdCastGroupManageReq {
    detail := &BrdCastGroupManageReq{
        GroupId:   string(bcg.Id),
        MemberIds: make([]string, len(bcg.Members)),
    }

    bcg.lock.Lock()
    for id, idx := range bcg.Members {
        detail.MemberIds[idx] = string(id)
    }
    bcg.lock.Unlock()

    return detail
}

func (bcg *BrdCastGroup) ResetMembers(memberIds... string) int32 {
    bcg.lock.Lock()

    for id := range bcg.Members {
        delete(bcg.Members, id)
    }

    for index, memberId := range memberIds{
        bcg.Members[IdString(memberId)] = int32(index)
        index += 1
    }
    bcg.lock.Unlock()

    return int32(len(memberIds))
}

func (bcg *BrdCastGroup) AddMember(memberId IdString) int32 {
    bcg.lock.Lock()

    index := len(bcg.Members)
    if _, ok := bcg.Members[IdString(memberId)]; !ok {
        bcg.Members[memberId] = int32(index)
        index += 1
    }

    bcg.lock.Unlock()

    return int32(index) //member num
}

func (bcg *BrdCastGroup) AddMembers(memberIds... string) int32 {
    bcg.lock.Lock()
    index := len(bcg.Members)
    for _, memberId := range memberIds{
        if _, ok := bcg.Members[IdString(memberId)]; !ok {
            bcg.Members[IdString(memberId)] = int32(index)
            index += 1
        }
    }

    bcg.lock.Unlock()

    return int32(index) //member num
}

func (bcg *BrdCastGroup) DelMember(memberId IdString) (int32, bool) {
    bcg.lock.Lock()

    index, ok := bcg.Members[memberId]
    if !ok {
        bcg.lock.Unlock()
        return 0, false
    }

    delete(bcg.Members, memberId)
    for id, idx := range bcg.Members {
        if idx > index {
            bcg.Members[id] = int32(idx) - 1
        }
    }
    num := len(bcg.Members)
    bcg.lock.Unlock()

    return int32(num), true //member num
}

func (bcg *BrdCastGroup) MemberNum() (int32) {
    bcg.lock.Lock()
    num := len(bcg.Members)
    bcg.lock.Unlock()

    return int32(num) //member num
}