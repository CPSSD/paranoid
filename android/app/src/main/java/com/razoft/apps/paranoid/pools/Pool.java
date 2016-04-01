package com.razoft.apps.paranoid.pools;

public class Pool {

    private String _fullname;
    private String _propername;
    private String _discovery;

    public Pool(String fullName, String properName){
        this._propername = properName;
        this._fullname = fullName;
        this._discovery = "paranoid.discovery.razoft.net:10101";
    }

    public Pool(String fullname, String propername, String discovery){
        this._propername = propername;
        this._fullname = fullname;
        this._discovery = discovery;
    }

    public void SetDiscovery(String discovery){
        this._discovery = discovery;
    }

    public String GetProperName(){
        return this._propername;
    }

    public String GetFullName(){
        return this._fullname;
    }

    public String GetDiscovery(){
        return this._discovery;
    }
}
