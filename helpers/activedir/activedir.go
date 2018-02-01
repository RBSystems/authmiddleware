package activedir

import (
	"fmt"
	"os"

	"github.com/mavricknz/ldap"
)

func GetGroupsForUser(userID string) ([]string, error) {
	groups := []string{}

	conn := ldap.NewLDAPConnection(
		"cad3.byu.edu",
		389)
	err := conn.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	search := ldap.NewSearchRequest(
		"ou=people,dc=byu,dc=local",
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		fmt.Sprintf("(Name=%s)", userID),
		[]string{"Name", "MemberOf"},
		nil,
	)
	username := os.Getenv("LDAP_USERNAME")
	password := os.Getenv("LDAP_PASSWORD")

	err = conn.Bind(username, password)
	if err != nil {
		panic(err)
	}

	res, err := conn.Search(search)
	if err != nil {
		panic(err)
	}

	//verify name
	for i := 0; i < len(res.Entries); i++ {
		name := res.Entries[i].GetAttributeValue("Name")
		if name != userID {
			continue
		}

		groupsTemp := res.Entries[0].GetAttributeValues("MemberOf")
		groups = translateGroups(groupsTemp)
	}

	//Extract groups
	fmt.Printf("results: %+v\n", res.Entries[0].GetAttributeValues("MemberOf")[0])

	return groups, nil
}

/*Format comes in 1, needs to be translated to 2.
  CN=AVS-TEC-SCCMAdmins,OU=Products,OU=AV Services,OU=OIT,DC=byu,DC=local
  byu.local/OIT/AV Services/Products/AVS-TEC-SCCMAdmins
*/
func translateGroups(groups []string) []string {
	toReturn := []string{}

	for _, entry := range groups {
		//TODO: See Comment
		toReturn = append(toReturn, entry)
	}
	return toReturn
}
