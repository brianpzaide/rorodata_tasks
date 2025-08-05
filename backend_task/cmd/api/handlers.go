package main

import (
	"errors"
	"net/http"
	"rorodata_backend_task/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) createCluster(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cluster, err := app.models.CreateCluster(input.Name, input.Region)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"cluster": cluster}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCluster(w http.ResponseWriter, r *http.Request) {

	cis := chi.URLParam(r, "cluster_id")
	clusterId, err := strconv.Atoi(cis)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("cluster_id must be a positive integer"))
		return
	}
	err = app.models.DeleteCluster(int64(clusterId))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrVMSExists):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "cluster successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createMachine(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string   `json:"name"`
		InstanceType string   `json:"instance_type"`
		Tags         []string `json:"tags"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cis := chi.URLParam(r, "cluster_id")
	clusterId, err := strconv.Atoi(cis)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("cluster_id must be a positive integer"))
		return
	}

	vm, err := app.models.CreateVM(int64(clusterId), input.Name, input.InstanceType, input.Tags)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, models.ErrNoClusterFound):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"VM": vm}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMachine(w http.ResponseWriter, r *http.Request) {
	cis := chi.URLParam(r, "cluster_id")
	clusterId, err := strconv.Atoi(cis)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("cluster_id must be a positive integer"))
		return
	}

	vmis := chi.URLParam(r, "vm_id")
	vmId, err := strconv.Atoi(vmis)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("vm_id must be a positive integer"))
		return
	}

	err = app.models.DeleteVM(int64(clusterId), int64(vmId))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, models.ErrNoClusterFound):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "machine successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) changeMachineState(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Operation string   `json:"operation"`
		Tags      []string `json:"tags"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	vms, err := app.models.Operate(input.Operation, input.Tags)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"vms_affected": vms}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
