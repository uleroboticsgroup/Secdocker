package docker

import (
	log "github.com/sirupsen/logrus"

	"niebla.unileon.es/DavidFerng/secdocker/config"
)

func panicForbbiden(itemForbidden string, itemType string) {
	log.Error("Forbidden ", itemType, ": ", itemForbidden)
}

// CheckIntersection compares the items between 2 slices and check if there is a match
func CheckIntersection(userItems []string, securityItems []string, itemType string) bool {
	check := false

	for _, item1 := range securityItems {
		for _, item2 := range userItems {
			if item1 == item2 {
				panicForbbiden(item1, itemType)
				check = true
			}
		}
	}

	return check
}

// CheckPermissions checks if container creation instructions pass all restrictions
func CheckPermissions(data ContainerOpts) bool {
	config := config.LoadConfig()
	check := true

	// Because of short-circuit evaluation
	if CheckIntersection(data.Ports, config.Restrictions.Ports, "port") {
		check = false
	}

	if CheckIntersection(data.Mounts, config.Restrictions.Mounts, "mount") {
		check = false
	}

	if CheckIntersection(data.Env, config.Restrictions.Environment, "environment") {
		check = false
	}

	if CheckIntersection([]string{data.User}, config.Restrictions.Users, "users") {
		check = false
	}

	if CheckIntersection([]string{data.Image}, config.Restrictions.Images, "images") {
		check = false
	}

	if data.Privileged == config.Restrictions.Privileged {
		log.Error("Forbbiden privileged")
		check = false
	}

	if CheckIntersection(data.SecurityPolicies, config.Restrictions.SecurityPolicies, "securityPolicies") {
		check = false
	}

	return check
}
