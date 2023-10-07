#!/bin/sh

DIR=$(dirname "$0")
INTERNAL_DIR="${DIR}"/../internal


mockgen -source="${INTERNAL_DIR}"/event/repository.go -destination="${INTERNAL_DIR}"/event/repository_mock.go -package=event -mock_names=Repository=MockRepository
mockgen -source="${INTERNAL_DIR}"/event/classifier/create_handler.go -destination="${INTERNAL_DIR}"/event/classifier/create_handler_mock.go -package=classifier -mock_names=CreateHandler=MockCreateHandler
mockgen -source="${INTERNAL_DIR}"/event/invitee/create_handler.go -destination="${INTERNAL_DIR}"/event/invitee/create_handler_mock.go -package=invitee -mock_names=CreateHandler=MockCreateHandler

