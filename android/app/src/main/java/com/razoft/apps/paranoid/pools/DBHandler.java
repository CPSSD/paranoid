package com.razoft.apps.paranoid.pools;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

import java.util.ArrayList;

public class DBHandler extends SQLiteOpenHelper {
    private static final int DATABASE_VERSION = 1;
    private static final String DATABASE_NAME = "paranoid.db";
    private static final String DATABASE_TABLE = "pools";

    private static final String COLUMN_ID = "_id";
    private static final String COLUMN_POOL = "poolname";
    private static final String COLUMN_DISCOVERY = "discovery";

    public DBHandler(Context context, String name, SQLiteDatabase.CursorFactory factory, int version ){
        super(context, DATABASE_NAME, factory, DATABASE_VERSION);
    }

    @Override
    public void onCreate(SQLiteDatabase db){
        String query = "CREATE TABLE " + DATABASE_TABLE + "(" +
                COLUMN_ID + " TEXT PRIMARY KEY," +
                COLUMN_POOL + " TEXT," +
                COLUMN_DISCOVERY + " TEXT" +
                ");";
        db.execSQL(query);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion){
        db.execSQL("DROP TABLE IF EXISTS" + DATABASE_TABLE);
        onCreate(db);
    }

    public boolean Add(Pool pool){

        ContentValues values = new ContentValues();
        values.put(COLUMN_ID, pool.GetProperName());
        values.put(COLUMN_POOL, pool.GetFullName());
        values.put(COLUMN_DISCOVERY, pool.GetDiscovery());
        SQLiteDatabase db = getWritableDatabase();
        try {
            db.insert(DATABASE_TABLE, null, values);
        } catch(Exception e){
            db.close();
            return false;
        }
        db.close();
        return true;
    }

    public Pool GetUsingProperName(String properName){
        SQLiteDatabase db = getReadableDatabase();
        Cursor c = db.rawQuery("SELECT * FROM " + DATABASE_TABLE + " WHERE _id=\"" + properName + "\"",null);
        if (!c.moveToFirst()){
            return null;
        }
        properName = c.getString(0);
        String fullName = c.getString(1);
        String discovery = c.getString(2);

        db.close();
        return new Pool(fullName, properName,discovery);
    }

    public ArrayList<Pool> GetAll(){
        ArrayList<Pool> pools = new ArrayList<>();

        SQLiteDatabase db = getReadableDatabase();
        Cursor cursor = db.rawQuery("SELECT * FROM " + DATABASE_TABLE + " WHERE 1=1", null);

        if(cursor.moveToFirst()) {
            while(!cursor.isAfterLast()){
                String properName = cursor.getString(0);
                String fullName = cursor.getString(1);
                String discovery = cursor.getString(2);
                pools.add(new Pool(fullName, properName, discovery));
                cursor.moveToNext();
            }
        }
        db.close();
        return pools;
    }
}
