package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/helpers"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestCreateSimpleReservation(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "csr_ok_"+id+"@test.com", "csr_ok_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource := helpers.CreateTestResource(t, "CSR Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		start := time.Now().Add(24 * time.Hour)
		end := start.Add(time.Hour)
		payload := map[string]interface{}{
			"ressource_id":   resource.ID.String(),
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handlers.CreateSimpleReservation(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var res models.Reservation
		json.Unmarshal(rec.Body.Bytes(), &res)
		if res.Status != "pending" {
			t.Errorf("Expected status 'pending', got '%s'", res.Status)
		}
		t.Cleanup(func() { config.DB.Delete(&res) })
	})

	t.Run("invalid_dates", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "csr_inv_"+id+"@test.com", "csr_inv_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource := helpers.CreateTestResource(t, "CSR Inv Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		start := time.Now().Add(48 * time.Hour)
		end := time.Now().Add(24 * time.Hour)
		payload := map[string]interface{}{
			"ressource_id":   resource.ID.String(),
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateSimpleReservation(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("missing_ressource_id", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "csr_nores_"+id+"@test.com", "csr_nores_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		start := time.Now().Add(24 * time.Hour)
		end := start.Add(time.Hour)
		payload := map[string]interface{}{
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateSimpleReservation(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("nonexistent_resource", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "csr_nrf_"+id+"@test.com", "csr_nrf_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		start := time.Now().Add(24 * time.Hour)
		end := start.Add(time.Hour)
		payload := map[string]interface{}{
			"ressource_id":   uuid.New().String(),
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateSimpleReservation(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("overlap", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "csr_ovlp_"+id+"@test.com", "csr_ovlp_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource := helpers.CreateTestResource(t, "CSR Overlap Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		start := time.Now().Add(72 * time.Hour)
		end := start.Add(time.Hour)

		// First reservation
		res1 := helpers.CreateTestReservation(t, user.ID, resource.ID)
		// Override times to our specific slot
		config.DB.Model(&res1).Updates(map[string]interface{}{
			"start_datetime": start,
			"end_datetime":   end,
		})
		t.Cleanup(func() { config.DB.Delete(&res1) })

		// Second overlapping reservation
		payload := map[string]interface{}{
			"ressource_id":   resource.ID.String(),
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateSimpleReservation(c)

		if rec.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusConflict, rec.Code, rec.Body.String())
		}
	})
}

func TestCreateGroupReservation(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "cgr_ok_"+id+"@test.com", "cgr_ok_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource1 := helpers.CreateTestResource(t, "CGR Resource1 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource1) })
		resource2 := helpers.CreateTestResource(t, "CGR Resource2 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource2) })

		start := time.Now().Add(24 * time.Hour)
		end := start.Add(time.Hour)
		payload := map[string]interface{}{
			"resource_ids":   []string{resource1.ID.String(), resource2.ID.String()},
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations/group", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handlers.CreateGroupReservation(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var res models.Reservation
		json.Unmarshal(rec.Body.Bytes(), &res)
		t.Cleanup(func() {
			config.DB.Where("reservation_id = ?", res.ID).Delete(&models.ReservationResource{})
			config.DB.Delete(&res)
		})
	})

	t.Run("too_few_resources", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "cgr_few_"+id+"@test.com", "cgr_few_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource := helpers.CreateTestResource(t, "CGR Few Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		start := time.Now().Add(24 * time.Hour)
		end := start.Add(time.Hour)
		payload := map[string]interface{}{
			"resource_ids":   []string{resource.ID.String()},
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations/group", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateGroupReservation(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("nonexistent_resource", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "cgr_nrf_"+id+"@test.com", "cgr_nrf_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource := helpers.CreateTestResource(t, "CGR NRF Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		start := time.Now().Add(24 * time.Hour)
		end := start.Add(time.Hour)
		payload := map[string]interface{}{
			"resource_ids":   []string{resource.ID.String(), uuid.New().String()},
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations/group", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateGroupReservation(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("overlap", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "cgr_ovlp_"+id+"@test.com", "cgr_ovlp_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource1 := helpers.CreateTestResource(t, "CGR Ovlp Resource1 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource1) })
		resource2 := helpers.CreateTestResource(t, "CGR Ovlp Resource2 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource2) })

		start := time.Now().Add(96 * time.Hour)
		end := start.Add(time.Hour)

		// Create overlapping existing reservation on resource1
		res1 := helpers.CreateTestReservation(t, user.ID, resource1.ID)
		config.DB.Model(&res1).Updates(map[string]interface{}{
			"start_datetime": start,
			"end_datetime":   end,
		})
		t.Cleanup(func() { config.DB.Delete(&res1) })

		// Now try to group reserve both
		payload := map[string]interface{}{
			"resource_ids":   []string{resource1.ID.String(), resource2.ID.String()},
			"start_datetime": start.Format(time.RFC3339),
			"end_datetime":   end.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/reservations/group", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.CreateGroupReservation(c)

		if rec.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusConflict, rec.Code, rec.Body.String())
		}
	})
}

func TestGetUserReservations(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "gur_ok_"+id+"@test.com", "gur_ok_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		resource := helpers.CreateTestResource(t, "GUR Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodGet, "/reservations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		err := handlers.GetUserReservations(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var reservations []models.Reservation
		json.Unmarshal(rec.Body.Bytes(), &reservations)
		if len(reservations) == 0 {
			t.Error("Expected at least one reservation")
		}
	})

	t.Run("missing_auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reservations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.GetUserReservations(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})
}

func TestGetAdminReservations(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "GAR Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "GAR Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "GAR Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "gar_"+id+"@test.com", "gar_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodGet, "/admin/reservations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())

		err := handlers.GetAdminReservations(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("no_company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/reservations", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.GetAdminReservations(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestApproveReservation(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "ApproveRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "ApproveRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveRes Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "approve_res_"+id+"@test.com", "approve_res_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+reservation.ID.String()+"/approve", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(reservation.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.ApproveReservation(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.Reservation
		config.DB.First(&updated, "id = ?", reservation.ID)
		if updated.Status != "approved" {
			t.Errorf("Expected status 'approved', got '%s'", updated.Status)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "ApproveRes NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+uuid.New().String()+"/approve", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.ApproveReservation(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "ApproveRes WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "ApproveRes WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "ApproveRes WC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveRes WC Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "approve_wc_"+id+"@test.com", "approve_wc_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+reservation.ID.String()+"/approve", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(reservation.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.ApproveReservation(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestRejectReservation(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RejectRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "RejectRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "RejectRes Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "reject_res_"+id+"@test.com", "reject_res_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+reservation.ID.String()+"/reject", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(reservation.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.RejectReservation(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.Reservation
		config.DB.First(&updated, "id = ?", reservation.ID)
		if updated.Status != "rejected" {
			t.Errorf("Expected status 'rejected', got '%s'", updated.Status)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RejectRes NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+uuid.New().String()+"/reject", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.RejectReservation(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "RejectRes WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "RejectRes WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "RejectRes WC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "RejectRes WC Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "reject_wc_"+id+"@test.com", "reject_wc_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+reservation.ID.String()+"/reject", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(reservation.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.RejectReservation(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestApproveGroupReservationResources(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "ApproveGR Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "ApproveGR Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource1 := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveGR Resource1 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource1) })
		resource2 := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveGR Resource2 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource2) })
		user := helpers.CreateTestUser(t, "approvegr_"+id+"@test.com", "approvegr_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		// Create a group reservation (no ressource_id)
		groupRes := models.Reservation{
			ID:            uuid.New(),
			UserID:        user.ID,
			StartDatetime: time.Now().Add(24 * time.Hour),
			EndDatetime:   time.Now().Add(25 * time.Hour),
			Status:        "pending",
		}
		config.DB.Create(&groupRes)
		t.Cleanup(func() { config.DB.Delete(&groupRes) })

		rr1 := models.ReservationResource{
			ID:            uuid.New(),
			ReservationID: groupRes.ID,
			ResourceID:    resource1.ID,
			Status:        "pending",
		}
		config.DB.Create(&rr1)
		t.Cleanup(func() { config.DB.Delete(&rr1) })

		rr2 := models.ReservationResource{
			ID:            uuid.New(),
			ReservationID: groupRes.ID,
			ResourceID:    resource2.ID,
			Status:        "pending",
		}
		config.DB.Create(&rr2)
		t.Cleanup(func() { config.DB.Delete(&rr2) })

		payload := map[string]interface{}{
			"resource_ids": []string{resource1.ID.String(), resource2.ID.String()},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+groupRes.ID.String()+"/resources/approve", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(groupRes.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.ApproveGroupReservationResources(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("not_a_group_reservation", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "ApproveGR Simple Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "ApproveGR Simple Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveGR Simple Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "approvegrsimple_"+id+"@test.com", "approvegrsimple_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		payload := map[string]interface{}{
			"resource_ids": []string{resource.ID.String()},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+reservation.ID.String()+"/resources/approve", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(reservation.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.ApproveGroupReservationResources(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "ApproveGR WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "ApproveGR WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "ApproveGR WC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource1 := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveGR WC Resource1 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource1) })
		resource2 := helpers.CreateTestResourceInPlace(t, place.ID, "ApproveGR WC Resource2 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource2) })
		user := helpers.CreateTestUser(t, "approvegrwc_"+id+"@test.com", "approvegrwc_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		groupRes := models.Reservation{
			ID:            uuid.New(),
			UserID:        user.ID,
			StartDatetime: time.Now().Add(24 * time.Hour),
			EndDatetime:   time.Now().Add(25 * time.Hour),
			Status:        "pending",
		}
		config.DB.Create(&groupRes)
		t.Cleanup(func() { config.DB.Delete(&groupRes) })

		rr1 := models.ReservationResource{
			ID:            uuid.New(),
			ReservationID: groupRes.ID,
			ResourceID:    resource1.ID,
			Status:        "pending",
		}
		config.DB.Create(&rr1)
		t.Cleanup(func() { config.DB.Delete(&rr1) })

		rr2 := models.ReservationResource{
			ID:            uuid.New(),
			ReservationID: groupRes.ID,
			ResourceID:    resource2.ID,
			Status:        "pending",
		}
		config.DB.Create(&rr2)
		t.Cleanup(func() { config.DB.Delete(&rr2) })

		payload := map[string]interface{}{
			"resource_ids": []string{resource1.ID.String(), resource2.ID.String()},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+groupRes.ID.String()+"/resources/approve", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(groupRes.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.ApproveGroupReservationResources(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestRejectGroupReservationResources(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RejectGR Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "RejectGR Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource1 := helpers.CreateTestResourceInPlace(t, place.ID, "RejectGR Resource1 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource1) })
		resource2 := helpers.CreateTestResourceInPlace(t, place.ID, "RejectGR Resource2 "+id)
		t.Cleanup(func() { config.DB.Delete(&resource2) })
		user := helpers.CreateTestUser(t, "rejectgr_"+id+"@test.com", "rejectgr_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		groupRes := models.Reservation{
			ID:            uuid.New(),
			UserID:        user.ID,
			StartDatetime: time.Now().Add(24 * time.Hour),
			EndDatetime:   time.Now().Add(25 * time.Hour),
			Status:        "pending",
		}
		config.DB.Create(&groupRes)
		t.Cleanup(func() { config.DB.Delete(&groupRes) })

		rr1 := models.ReservationResource{
			ID:            uuid.New(),
			ReservationID: groupRes.ID,
			ResourceID:    resource1.ID,
			Status:        "pending",
		}
		config.DB.Create(&rr1)
		t.Cleanup(func() { config.DB.Delete(&rr1) })

		rr2 := models.ReservationResource{
			ID:            uuid.New(),
			ReservationID: groupRes.ID,
			ResourceID:    resource2.ID,
			Status:        "pending",
		}
		config.DB.Create(&rr2)
		t.Cleanup(func() { config.DB.Delete(&rr2) })

		payload := map[string]interface{}{
			"resource_ids": []string{resource1.ID.String(), resource2.ID.String()},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+groupRes.ID.String()+"/resources/reject", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(groupRes.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.RejectGroupReservationResources(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("not_a_group_reservation", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RejectGR Simple Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "RejectGR Simple Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "RejectGR Simple Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "rejectgrsimple_"+id+"@test.com", "rejectgrsimple_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		payload := map[string]interface{}{
			"resource_ids": []string{resource.ID.String()},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/reservations/"+reservation.ID.String()+"/resources/reject", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(reservation.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.RejectGroupReservationResources(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
