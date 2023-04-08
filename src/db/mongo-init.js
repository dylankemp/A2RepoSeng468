db = db.getSiblingDB("admin") 
db.auth("admin","admin") 
db = db.getSiblingDB("social_app") 

db.createUser({ 
    user: "user", 
    pwd: "user", 
    roles: [ 
        { 
            role: "readWrite", 
            db: "social_app"
        }
    ]
}); 

db.createCollection("users")