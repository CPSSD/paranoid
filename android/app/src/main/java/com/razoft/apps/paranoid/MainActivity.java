package com.razoft.apps.paranoid;

import android.app.AlertDialog;
import android.app.Dialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.graphics.Color;
import android.os.Bundle;
import android.support.design.widget.FloatingActionButton;
import android.support.design.widget.NavigationView;
import android.support.design.widget.Snackbar;
import android.support.v4.view.GravityCompat;
import android.support.v4.widget.DrawerLayout;
import android.support.v7.app.ActionBar;
import android.support.v7.app.ActionBarDrawerToggle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.view.menu.MenuBuilder;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.SubMenu;
import android.view.View;
import android.support.v7.widget.Toolbar;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

import com.razoft.apps.paranoid.pools.*;

import java.util.ArrayList;

public class MainActivity extends AppCompatActivity
    implements NavigationView.OnNavigationItemSelectedListener {

    private boolean pfsdOnline;
    private DBHandler dbHandler;

    private static final String TAG = "paranoid";

    Pool current;

    @Override
    public void onCreate(Bundle savedInstanceState){
        super.onCreate(savedInstanceState);

        // Create the datebase handler
        dbHandler = new DBHandler(this,null,null,1);

        //Create the pool dialog
        newPool();
        fullnameText = (EditText) findViewById(R.id.newpool_full);
        propernameText = (EditText) findViewById(R.id.newpool_proper);
        discoveryText = (EditText) findViewById(R.id.newpool_discovery);

        // Load the activity
        setContentView(R.layout.activity_main);

        // Load the toolbar and activate the drawer
        Toolbar toolbar = (Toolbar) findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);

        DrawerLayout drawer = (DrawerLayout) findViewById(R.id.drawer_layout);
        ActionBarDrawerToggle toggle = new ActionBarDrawerToggle(
                this, drawer, toolbar,
                R.string.sidebar_drawer_open,
                R.string.sidebar_drawer_close);
        drawer.setDrawerListener(toggle);
        toggle.syncState();

        NavigationView navView = (NavigationView) findViewById(R.id.nav_view);

        dbHandler.Add(new Pool("Hello2", "world2"));
        dbHandler.Add(new Pool("Hello3", "world3"));

        ArrayList<Pool> pools = dbHandler.GetAll();

        navView.setNavigationItemSelectedListener(this);

        Menu list = (Menu) navView.getMenu();

        for(Pool p : pools) {
            list.add(p.GetProperName());
        }


        // Create the Floating Action Button
        createFab();
    }

    @Override
    public void onBackPressed(){
        DrawerLayout drawer = (DrawerLayout) findViewById(R.id.drawer_layout);
        if(drawer.isDrawerOpen(GravityCompat.START)){
            drawer.closeDrawer(GravityCompat.START);
        } else {
            super.onBackPressed();
        }
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu){
        // Populate the options
        getMenuInflater().inflate(R.menu.options, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item){
        int id = item.getItemId();

        Intent intent = null;

        switch(id){
            case R.id.action_about:
                intent = new Intent(this, AboutActivity.class);
                break;
            case R.id.action_settings:
                intent = new Intent(this, SettingsActivity.class);
                break;
        }

        startActivity(intent);
        return super.onOptionsItemSelected(item);
    }

    @SuppressWarnings("StatementWithEmptyBody")
    @Override
    public boolean onNavigationItemSelected(MenuItem item){
        int id = item.getItemId();

        switch(id){
            case R.id.pool_properties:
                break;
            case R.id.menu_add_pool:
                dialog.show();
                break;
            default:
                String proper = item.getTitle().toString();
                switchToPool(proper);
                break;
        }

        DrawerLayout drawer = (DrawerLayout) findViewById(R.id.drawer_layout);
        drawer.closeDrawer(GravityCompat.START);
        return true;
    }

    // Creates a fab at the bottom of the page
    public void createFab(){
        final FloatingActionButton fab = (FloatingActionButton) findViewById(R.id.fab);

        if(pfsdOnline){
            fab.setImageResource(R.drawable.ic_media_stop);
            updateContent();
        }

        fab.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View v) {
                String message = "";

                if (pfsdOnline) {
                    message = getString(R.string.status_pfsd_stop);
                    pfsdOnline = false;
                    stopPfsd();
                    fab.setImageResource(R.drawable.ic_media_play);
                } else {
                    if (startPfsd()) {
                        pfsdOnline = true;
                        message = getString(R.string.status_pfsd_success);
                        fab.setImageResource(R.drawable.ic_media_stop);
                    } else {
                        pfsdOnline = false;
                    }
                }

                Snackbar.make(v, message, Snackbar.LENGTH_LONG)
                        .setAction(android.R.string.ok, null).show();

                updateContent();
            }
        });
    }

    public boolean startPfsd(){
        //TODO: Write code to start pfsd here

        return true;
    }

    public boolean stopPfsd(){
        //TODO: Write code to stop pfsd
        pfsdOnline = false;
        return true;
    }

    public NodeInfo getStatusInfo(){
        NodeInfo ni= new NodeInfo();
        ni.isOnline = true;
        ni.localAddress = "127.0.0.1:12345";
        ni.peers = new String[3];
        ni.peers[0] = "10.0.0.1:23456";
        ni.peers[1] = "10.0.0.2:34567";
        ni.peers[2] = "10.0.0.3:45678";

        return ni;
    }

    private MainActivity getContext(){
        return this;
    }

    private void updateContent(){
        TextView connectedPeers = (TextView) findViewById(R.id.status_connected_peers);
        TextView currentStatus = (TextView) findViewById(R.id.status_current);
        TextView localAddress = (TextView) findViewById(R.id.status_local_address);
        if(pfsdOnline){
            ((ImageView)findViewById(R.id.pfsd_status)).setBackgroundColor(Color.GREEN);

            String allPeers = "";
            for(String s :getStatusInfo().peers){
                allPeers += s+"\n";
            }
            connectedPeers.setText(allPeers);
            currentStatus.setText("Connected");
            localAddress.setText(getStatusInfo().localAddress);
        } else {
            ((ImageView)findViewById(R.id.pfsd_status)).setBackgroundColor(Color.RED);
            connectedPeers.setText("None");
            currentStatus.setText("Disconnected");
            localAddress.setText("");
        }
    }

    private EditText fullnameText;
    private EditText propernameText;
    private EditText discoveryText;
    private Dialog dialog;

    public void newPool(){
        final AlertDialog.Builder newPoolDialogBuilder = new AlertDialog.Builder(this);

        LayoutInflater inflater = LayoutInflater.from(this);
        newPoolDialogBuilder.setTitle(R.string.newpool_title);
        newPoolDialogBuilder.setView(inflater.inflate(R.layout.fragment_new_pool_dialog, null));

        newPoolDialogBuilder.setPositiveButton(R.string.newpool_add,
                new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {

                        String fullname = fullnameText.toString();
                        String propername = propernameText.toString();
                        String disovery = discoveryText.toString();

                        dbHandler.Add(new Pool(fullname, propername, disovery));
                        Toast.makeText(getContext(), R.string.newpool_add_success + fullname,
                                Toast.LENGTH_LONG).show();
                        supportInvalidateOptionsMenu();
                    }
                }
        );

        newPoolDialogBuilder.setNegativeButton(R.string.newpool_cancel,
                new DialogInterface.OnClickListener(){
                    @Override
                public void onClick(DialogInterface dialog, int which){
                        dialog.dismiss();
                    }
                });

        dialog = newPoolDialogBuilder.create();
    }

    public void switchToPool(String propername){
        current = dbHandler.GetUsingProperName(propername);

        ((TextView)findViewById(R.id.sidebar_header_pool_full)).setText(current.GetFullName());
        ((TextView)findViewById(R.id.sidebar_header_pool_proper)).setText(current.GetProperName());

//        propernameText.setText(current.GetProperName());
//        fullnameText.setText(current.GetFullName());
//        discoveryText.setText(current.GetDiscovery());

        stopPfsd();
        updateContent();
        supportInvalidateOptionsMenu();

    }
}

class NodeInfo {
    public boolean isOnline;
    public String localAddress;
    public String []peers;
}