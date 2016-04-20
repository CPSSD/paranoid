class DB {
  constructor(name, ver){
    console.info("Database initialized");
    this.m_version = ver;
    this.m_dbname = name;
  }

  get(){

  }

  set(){

  }

  version(){
    return this.m_version;
  }

  onUpgrade(oldVersion, newVersion){

  }
}

const TemporaryDB = 1;
const PermanentDB = 2;

class KeyValueDB extends DB {
  constructor(name, version, databasePersistance){
    super(name, version, databasePersistance);
    this.m_name = name;

    if(databasePersistance == TemporaryDB){
      this.m_db = sessionStorage;
    }
    if(databasePersistance == PermanentDB){
      this.m_db = localStorage;
    }

    this.m_db[this.m_name+"__version__"] = version;
  }

  get(key){
    return this.m_db[this.m_name+key];
  }

  set(key, value){
    this.m_db[this.m_name+key] = value;
  }
}

//TODO: Implement RelationalDB
class RelationalDB extends DB {

}
