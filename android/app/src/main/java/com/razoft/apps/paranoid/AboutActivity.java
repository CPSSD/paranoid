package com.razoft.apps.paranoid;

import android.content.Intent;
import android.net.Uri;
import android.support.v7.app.AppCompatActivity;
import android.view.Gravity;
import android.view.View;
import android.widget.LinearLayout;
import android.os.Bundle;
import android.widget.TextView;

public class AboutActivity extends AppCompatActivity{

    @Override
    protected void onCreate(Bundle savedInstanceState){
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_about);

        // Populate the about page
        LinearLayout androidDevs = (LinearLayout) findViewById(R.id.about_android_devs);
        LinearLayout paranoidDevs = (LinearLayout) findViewById(R.id.about_paranoid_devs);
        LinearLayout projectMentor = (LinearLayout) findViewById(R.id.about_project_mentor);

        TextView wojciech = new TextView(this);
        wojciech.setText("Wojciech Bednarzak");
        wojciech.setGravity(Gravity.CENTER_HORIZONTAL);
        wojciech.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse("https://github.com/VoyTechnology")));
            }
        });
        androidDevs.addView(wojciech);



        TextView terry = new TextView(this);
        terry.setGravity(Gravity.CENTER_HORIZONTAL);
        terry.setText("Terry Bolt");
        terry.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse("https://github.com/GoldenBadger")));
            }
        });
        paranoidDevs.addView(terry);

        TextView conor = new TextView(this);
        conor.setGravity(Gravity.CENTER_HORIZONTAL);
        conor.setText("Conor Griffin");
        conor.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse("https://github.com/ConorGriffin37")));
            }
        });
        paranoidDevs.addView(conor);

        TextView sean = new TextView(this);
        sean.setGravity(Gravity.CENTER_HORIZONTAL);
        sean.setText("Sean Healy");
        sean.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse("https://github.com/SeanHealy33")));
            }
        });
        paranoidDevs.addView(sean);

        TextView mladen = new TextView(this);
        mladen.setGravity(Gravity.CENTER_HORIZONTAL);
        mladen.setText("Mladen Kajic");
        mladen.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse("https://github.com/CanOpener")));
            }
        });
        paranoidDevs.addView(mladen);

        TextView stephen = new TextView(this);
        stephen.setGravity(Gravity.CENTER_HORIZONTAL);
        stephen.setText("Stephen Blott");
        stephen.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse("https://github.com/smblott-github")));
            }
        });
        projectMentor.addView(stephen);

        // Show Version
        String version = BuildConfig.VERSION_NAME;
        TextView versionElement = null;
        versionElement = (TextView) findViewById(R.id.about_version);
        versionElement.setText("v"+version);
    }
}
