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

import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:equatable/equatable.dart';

import 'mobileapi/mobileapi.dart';
import 'topic_bloc.dart';
import 'topic_map_view_models.dart';

class LibraryBloc extends Bloc<LibraryEvent, LibraryState> {
  final QueryApi queryApi;
  final CommandApi commandApi;

  LibraryBloc({
    @required this.queryApi,
    @required this.commandApi,
  }) {
    print("LibraryBloc created: ${this}");
  }

  @override
  LibraryState get initialState {
    print("creating initial state");
    return LibraryState();
  }

  @override
  Stream<LibraryState> mapEventToState(LibraryEvent event) async* {
    if (event is LibraryAppStartedEvent) {
      yield LibraryState(loading: true);
      yield await _loadTopicMaps();
    }

    if (event is LibraryReloadEvent) {
      yield await _loadTopicMaps();
    }

    if (event is LibraryTopicMapDeletedEvent) {
      DeleteTopicMapRequest request = DeleteTopicMapRequest();
      request.topicMapId = event.topicMapId;
      await commandApi.deleteTopicMap(request);
      yield await _loadTopicMaps();
    }
  }

  Future<LibraryState> _loadTopicMaps() async {
    LibraryState next;
    await queryApi.getTopicMaps(GetTopicMapsRequest()).then((response) {
      next = LibraryState(
        topicMaps: (response.topicMaps ?? const [])
            .map((tm) => TopicMapViewModel(tm))
            .toList(growable: false),
      );
    }).catchError((error) {
      next = LibraryState(error: error.toString());
    });
    return next;
  }

  TopicBloc createTopicBloc({TopicViewModel viewModel}) => TopicBloc(
        queryApi: queryApi,
        commandApi: commandApi,
        libraryBloc: this,
        viewModel: viewModel,
      );
}

class LibraryState {
  final List<TopicMapViewModel> topicMaps;
  final bool loading;
  final String error;

  LibraryState({
    this.topicMaps = const [],
    this.loading = false,
    this.error,
  }) : assert(topicMaps != null);

  @override
  String toString() =>
      "LibraryState(${topicMaps.length}, ${loading}, ${error})";
}

class LibraryEvent extends Equatable {}

class LibraryAppStartedEvent extends LibraryEvent {}

class LibraryReloadEvent extends LibraryEvent {}

class LibraryTopicMapDeletedEvent extends LibraryEvent {
  final Int64 topicMapId;

  LibraryTopicMapDeletedEvent(this.topicMapId)
      : assert(topicMapId != null && topicMapId != 0);
}
