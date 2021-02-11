package operators

func GetNodes() {
	kube := pipe.retrieveKube()

	for {
		msg := <-pipe.getChan
		node := msg.Node

		dep, err := kube.Get(node.ToDeployment(), metav1.GetOptions{})
		code := treatK8sError(err, cerrors.InsprErrorCode(6))
		if code != 0 {
			insprError := generateInsprError(5, err, "error in node back end update")
			pipe.responseChan <- generateInsprResponseMessge(insprError, msg.Spec)
			continue
		}

		node.NodeStatus, err = pipe.getNodeStatus(dep)
		fmt.Println(node.NodeStatus)
		if err != nil {
			logrus.Println("unable to get pod status")
			insprError := generateInsprError(5, err, "unable to get pod status")
			pipe.responseChan <- generateInsprResponseMessge(insprError, msg.Spec)
		}
		_, err = pipe.registry.PutNode(node)
		if err != nil {
			insprError := generateInsprError(5, err, "error in node back end update")
			pipe.responseChan <- generateInsprResponseMessge(insprError, msg.Spec)
		} else {
			pipe.responseChan <- generateInsprResponseMessge(nil, msg.Spec)
		}
	}
}
