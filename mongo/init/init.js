db.createUser({
  user: "totallyNotAdmin",
  pwd: "secret1337",
  roles: [{role: "readWrite", db: "pasteBinItems"}]
})
