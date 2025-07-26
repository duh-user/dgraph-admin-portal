package user

// all queries need to start with the name "query" to work with our query handler
const (
	QFIELDSUSER = `
		uid
		name
		user_name
		email
		role {
			uid
			role_name
		}
		pass_hash
		date_created
		last_modified
		last_seen
	`
	QBYNAMEEXACT = `
		query query($name: string) {
			query(func: eq(name, $name)) {
				` + QFIELDSUSER + `
			}	
		}`

	QBYNAMEFUZZY = `
		query query($name: string) {
			query(func: match(name, $name, 25))	{
				` + QFIELDSUSER + `
			}	
		}`

	QBYEMAILEXACT = `
		query query($email: string) {
			query(func: eq(email, $email)) {
				` + QFIELDSUSER + `
			}	
		}`

	QBYEMAILFUZZY = `
		query query($email: string) {
			query(func: match(email, $email, 25)) {
				` + QFIELDSUSER + `
			}	
		}`

	QBYUNAMEEXACT = `
		query query($user_name: string) {
			query(func: eq(user_name, $user_name)) {
				` + QFIELDSUSER + `
			}	
		}`

	QBYUNAMEFUZZY = `
		query query($user_name: string) {
			query(func: match(user_name, $user_name, 25)) {
				` + QFIELDSUSER + `
			}	
		}`

	QBYUID = `
	query query($uid: string) {
		query(func: uid($uid)) {
			` + QFIELDSUSER + `
			}	
		}`

	QBYROLE = `
		query query($role: string) {
			query(func: eq(role, $role)) {
				` + QFIELDSUSER + `
			}
		}`

	QALLUSERS = `
		query query($role: string) {
			query(func: eq(role, $role)) {
				` + QFIELDSUSER + `
			}
		}`
)
