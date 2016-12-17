package com

import (
    . "types"
    "sync"
)

var brdCastGroupsByServer = make(map[IdString]*BrdCastGroupManage, 0)

func FindBrdCastGroup(serverId IdString, groupId IdString) (*BrdCastGroup, bool) {
    bm, ok := brdCastGroupsByServer[serverId]
    if !ok {
        brdCastGroupsByServer[serverId] = &BrdCastGroupManage{
            groups: make(map[IdString]*BrdCastGroup),
        }
        return nil, false
    }

    grp, ok := bm.FindGroup(groupId)
    if !ok {
        return nil, false
    }

    return grp, true
}

func DestroyBrdCastGroup(serverId IdString, groupId IdString) (*BrdCastGroup, bool) {
    bm, ok := brdCastGroupsByServer[serverId]
    if !ok {
        brdCastGroupsByServer[serverId] = &BrdCastGroupManage{
            groups: make(map[IdString]*BrdCastGroup),
        }
        return nil, false
    }

    grp, ok := bm.DestroyGroup(groupId)
    if !ok {
        return nil, false
    }

    return grp, true
}

func NewBrdCastGroup4Server(serverId IdString, groupId IdString) (*BrdCastGroup) {
    bm, ok := brdCastGroupsByServer[serverId]
    if !ok {
        brdCastGroupsByServer[serverId] = &BrdCastGroupManage{
            groups: make(map[IdString]*BrdCastGroup),
        }
    }

    grp, ok := bm.FindGroup(groupId)
    if !ok {
        grp = bm.NewGroup(groupId)
    }

    return grp
}


type BrdCastGroupManage struct {
    groups map[IdString]*BrdCastGroup
    lock  sync.Mutex
}

func (bm *BrdCastGroupManage) NewGroup(creator IdString) *BrdCastGroup {
    bm.lock.Lock()
    if grp, ok := bm.groups[creator]; ok {
        bm.lock.Unlock()
        return grp
    }

    grp := NewBrdCastGroup(creator)
    bm.groups[grp.Id] = grp
    bm.lock.Unlock()

    grp.AddMember(creator)

    return grp
}

func (bm *BrdCastGroupManage) FindGroup(groupId IdString) (*BrdCastGroup, bool) {
    bm.lock.Lock()
    grp, ok := bm.groups[groupId]
    bm.lock.Unlock()

    return grp, ok
}

func (bm *BrdCastGroupManage) DestroyGroup(groupId IdString) (*BrdCastGroup, bool) {
    bm.lock.Lock()
    grp, ok := bm.groups[groupId]
    if ok {
        delete(bm.groups, groupId)
    }
    bm.lock.Unlock()

    return grp, ok
}