package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
	"bitbucket.org/dragoninfo/temtum-go-master/internal/operations"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

//change leader operation
func (s *Service) ChangeActiveLieder() {
	{
		log.Info("WRITE ROLE")
		role := isTransactionReader(&s.isTransactionReader)

		err := s.serv.SendToChanel("leaderHash", operations.SendRoleOperation{Msg: role})
		if err != nil {
			log.Errorf("role message not send to old leader", err)
		}
		//err = s.serv.SendToChanel(SelectNextLeader(s.nodeList), operations.SendRoleOperation{Msg: role})
		//if err != nil {
		//	log.Errorf("role message not send to new leader", err)
		//}
	}
}

func isTransactionReader(isReader *bool) string {
	role := activeRole
	if *isReader {
		role = passiveRole
		*isReader = false
	} else {
		*isReader = true
	}
	return role
}

//find active leader from node list
func findActiveLeader(nList []data.Status) string {
	for _, v := range nList {
		if v.Stat == ActiveLeader {
			return v.Hash
		}
	}
	return ""
}

//choose a leader who will become active
func selectNextLeader(nList []data.Status) string {
	peers := make([]data.Status, 0)

	for _, v := range nList {
		if v.Stat == Enable {
			peers = append(peers, v)
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	peerIndex := rand.Intn(len(peers))

	return nList[peerIndex].Hash
}
