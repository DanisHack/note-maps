// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'mobileapi.dart';

class TopicPage extends StatelessWidget {
  TopicPage({Key key, this.title}) : super(key: key);

  final String title;

  Widget noteTile(BuildContext context) {
    return Card(
      child: Placeholder(
        fallbackHeight: 80,
      ),
    );
  }

  Widget roleTile(BuildContext context) {
    return ListTile(
      leading: FlutterLogo(),
      title: Placeholder(
        fallbackHeight: 20,
      ),
      trailing: Icon(Icons.more_vert),
      onTap: () {
        Navigator.push(
          context,
          MaterialPageRoute(builder: (context) => TopicPage(title: "Topic")),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: <Widget>[
          SliverAppBar(
            pinned: true,
            snap: false,
            floating: false,
            expandedHeight: 160.0,
            flexibleSpace: FlexibleSpaceBar(
              title: Text(title),
              //background: Image.asset(..., fit: BoxFit.fill)
            ),
          ),
          SliverPadding(
            padding: const EdgeInsets.all(8.0),
            sliver: SliverList(
              delegate: SliverChildListDelegate(<Widget>[
                noteTile(context),
                noteTile(context),
                Divider(),
                roleTile(context),
                roleTile(context),
                roleTile(context),
              ]),
            ),
          ),
        ],
      ),
      bottomNavigationBar: BottomAppBar(
        child: Container(
          height: 50.0,
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
        tooltip: 'Create a Note',
        child: Icon(Icons.add),
      ),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
    );
  }
}
