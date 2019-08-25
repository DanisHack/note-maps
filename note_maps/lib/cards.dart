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
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import 'common_widgets.dart';
import 'mobileapi/controllers.dart';

class _FutureTextField extends StatelessWidget {
  final Future<TextEditingController> futureTextController;
  final bool autofocus;
  final TextCapitalization textCapitalization;
  final TextStyle style;

  _FutureTextField(
    this.futureTextController, {
    this.autofocus = false,
    this.textCapitalization,
    this.style,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<TextEditingController>(
      future: futureTextController,
      initialData: null,
      builder: (_, snapshot) {
        switch (snapshot.connectionState) {
          case ConnectionState.done:
            if (snapshot.hasError) {
              return ErrorIndicator();
            }
            return TextField(
              controller: snapshot.data,
              textCapitalization: textCapitalization,
              autofocus: autofocus,
              style: style,
              decoration: InputDecoration(border: InputBorder.none),
            );
          default:
            return CircularProgressIndicator();
        }
      },
    );
  }
}

class NameCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    var controller = Provider.of<NameController>(context);
    return ValueListenableBuilder<NameState>(
      valueListenable: controller,
      builder: (context, nameState, _) => Card(
        child: Row(
          children: <Widget>[
            Container(width: 48),
            Expanded(
              child: _FutureTextField(
                controller.valueTextController,
                textCapitalization: TextCapitalization.words,
                autofocus: true,
                style: Theme.of(context).textTheme.headline,
              ),
            ),
            _noteMenuButton(context, controller),
          ],
        ),
      ),
    );
  }
}

class OccurrenceCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    var controller = Provider.of<OccurrenceController>(context);
    return ValueListenableBuilder(
      valueListenable: controller,
      builder: (context, occurrenceState, _) => Card(
        child: Row(
          children: <Widget>[
            Container(width: 48),
            Expanded(
              child: _FutureTextField(
                controller.valueTextController,
                textCapitalization: TextCapitalization.sentences,
                style: Theme.of(context).textTheme.body2,
              ),
            ),
            _noteMenuButton(context, controller),
          ],
        ),
      ),
    );
  }
}

Widget _noteMenuButton(BuildContext context, NoteMapItemController controller) {
  return PopupMenuButton<NoteOption>(
    onSelected: (NoteOption choice) {
      switch (choice) {
        case NoteOption.delete:
          controller.delete().catchError(( error) {
            Scaffold.of(context)
                .showSnackBar(SnackBar(content: Text(error.toString())));
            return null;
          });
          break;
      }
    },
    itemBuilder: (BuildContext context) => <PopupMenuEntry<NoteOption>>[
      const PopupMenuItem<NoteOption>(
        value: NoteOption.delete,
        child: Text('Delete'),
      ),
    ],
  );
}

enum NoteOption { delete }
enum RoleOption {
  editRole,
  editAssociation,
}
