package main

type MapAOIClient struct {

}

func (pthis *MapAOIClient)Ntf_node_in_view(nodeID int, nodeIDInView int) {

}

func (pthis *MapAOIClient)Ntf_node_out_view(nodeID int, nodeIDOutView int) {

}

func ConstructMapAOIClient() *MapAOIClient {
	aoi := &MapAOIClient{}

	return aoi
}