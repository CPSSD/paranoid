class DB {
  constructor(ver){
    console.info("Database initialized");
    this.v = ver;
  }

  get(){

  }

  set(){

  }

  version(){
    return this.v;
  }

  onUpgrade(oldVersion, newVersion){

  }
}

const TemporaryDB = 1;
const PermanentDB = 2;

class KeyValueDB extends DB {
  constructor(ver, databasePersistance){
    super(ver, databasePersistance);

    if(databasePersistance == TemporaryDB){
      this.db = sessionStorage;
    }
    if(databasePersistance == PermanentDB){
      this.db = localStorage;
    }

    this.db["__version__"] = ver;
  }

  get(key){
    return this.db[key];
  }

  set(key, value){
    this.db[key] = value;
  }
}

//TODO: Implement RelationalDB
class RelationalDB extends DB {

}
