SELECT tag 
FROM channel
WHERE tag = $1 OR name = $1
