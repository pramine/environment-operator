package web

// DeployRequest represents POST request body to perform deployments.
//  * Name of the service to update
//  * Application image part (full construct from util.DockerImage )
//  * Version application version
type DeployRequest struct {
	Name        string `json:"name"`
	Application string `json:"application,omitempty"`
	Version     string
}
