package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func PubspecDependency() string {
	return "  flutter_bloc: ^8.1.6\n  equatable: ^2.0.5"
}

func GenerateBlocFiles(projectPath, featureName string) error {
	pascal := toPascalCase(featureName)
	base := filepath.Join(projectPath, fmt.Sprintf("lib/features/%s/presentation/bloc", featureName))

	files := map[string]string{
		filepath.Join(base, fmt.Sprintf("%s_bloc.dart", featureName)):  blocDart(pascal, featureName),
		filepath.Join(base, fmt.Sprintf("%s_event.dart", featureName)): eventDart(pascal),
		filepath.Join(base, fmt.Sprintf("%s_state.dart", featureName)): stateDart(pascal),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	result := ""
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

func blocDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '%s_event.dart';
import '%s_state.dart';

class %sBloc extends Bloc<%sEvent, %sState> {
  %sBloc() : super(%sInitial()) {
    on<%sLoadRequested>(_onLoadRequested);
  }

  Future<void> _onLoadRequested(
    %sLoadRequested event,
    Emitter<%sState> emit,
  ) async {
    emit(%sLoading());
    try {
      // TODO: implement load logic
      emit(const %sLoaded(data: []));
    } catch (e) {
      emit(%sError(message: e.toString()));
    }
  }
}
`, feature, feature, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}

func eventDart(pascal string) string {
	return fmt.Sprintf(`import 'package:equatable/equatable.dart';

abstract class %sEvent extends Equatable {
  const %sEvent();

  @override
  List<Object?> get props => [];
}

class %sLoadRequested extends %sEvent {
  const %sLoadRequested();
}
`, pascal, pascal, pascal, pascal, pascal)
}

func stateDart(pascal string) string {
	return fmt.Sprintf(`import 'package:equatable/equatable.dart';

abstract class %sState extends Equatable {
  const %sState();

  @override
  List<Object?> get props => [];
}

class %sInitial extends %sState {}

class %sLoading extends %sState {}

class %sLoaded extends %sState {
  final List<dynamic> data;
  const %sLoaded({required this.data});

  @override
  List<Object?> get props => [data];
}

class %sError extends %sState {
  final String message;
  const %sError({required this.message});

  @override
  List<Object?> get props => [message];
}
`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}
