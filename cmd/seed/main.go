package main

import (
	"fmt"
	"log"
	"time"

	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func hash(password string) []byte {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt error: %v", err)
	}
	return h
}

type seedResource struct {
	name     string
	category string
	status   string
}

type seedUser struct {
	email    string
	username string
	phone    int
	role     string
}

type seedPlace struct {
	name        string
	description string
	resources   []seedResource
}

type seedCompany struct {
	name   string
	places []seedPlace
	admin  seedUser
	users  []seedUser
}

func main() {
	if godotenv.Load(".env") != nil {
		godotenv.Load("../../.env")
	}
	config.ConnectDatabase()
	db := config.DB

	fmt.Println("🌱 Démarrage du seed...")

	// =====================
	// CATEGORIES
	// =====================
	categories := []models.Category{
		{ID: uuid.New(), Name: "Salle de réunion"},
		{ID: uuid.New(), Name: "Bureau individuel"},
		{ID: uuid.New(), Name: "Espace coworking"},
		{ID: uuid.New(), Name: "Salle de conférence"},
		{ID: uuid.New(), Name: "Laboratoire"},
		{ID: uuid.New(), Name: "PC"},
		{ID: uuid.New(), Name: "Projecteur"},
		{ID: uuid.New(), Name: "Imprimante"},
	}
	catMap := map[string]uuid.UUID{}
	for _, cat := range categories {
		var existing models.Category
		if db.Where("name = ?", cat.Name).First(&existing).Error != nil {
			db.Create(&cat)
			catMap[cat.Name] = cat.ID
			fmt.Printf("  ✅ Catégorie : %s\n", cat.Name)
		} else {
			catMap[cat.Name] = existing.ID
			fmt.Printf("  ⏭️  Catégorie déjà existante : %s\n", cat.Name)
		}
	}

	// =====================
	// ENTREPRISES + LIEUX + ADMINS + USERS + RESSOURCES
	// =====================
	companies := []seedCompany{
		{
			name: "TechCorp Solutions",
			admin: seedUser{"admin.techcorp@spacebook.fr", "admin_techcorp", 600100200, "admin"},
			users: []seedUser{
				{"alice.martin@techcorp.fr", "alice_martin", 600111001, "user"},
				{"bob.dupont@techcorp.fr", "bob_dupont", 600111002, "user"},
				{"claire.lefebvre@techcorp.fr", "claire_lefebvre", 600111003, "user"},
				{"david.moreau@techcorp.fr", "david_moreau", 600111004, "user"},
			},
			places: []seedPlace{
				{
					name:        "Siège Paris – Tour Montparnasse",
					description: "Bureaux principaux au cœur de Paris",
					resources: []seedResource{
						{"Salle Neptune", "Salle de réunion", "available"},
						{"Salle Saturne", "Salle de réunion", "available"},
						{"Grande Salle Prestige", "Salle de conférence", "available"},
						{"Bureau 301", "Bureau individuel", "available"},
						{"Bureau 302", "Bureau individuel", "available"},
						{"Espace Flex A", "Espace coworking", "available"},
						{"PC Paris 01", "PC", "available"},
						{"PC Paris 02", "PC", "available"},
						{"PC Paris 03", "PC", "unavailable"},
						{"Projecteur Paris 01", "Projecteur", "available"},
						{"Projecteur Paris 02", "Projecteur", "available"},
						{"Imprimante Paris 01", "Imprimante", "available"},
						{"Imprimante Paris 02", "Imprimante", "unavailable"},
					},
				},
				{
					name:        "Bureaux Lyon – Part-Dieu",
					description: "Antenne régionale Lyon",
					resources: []seedResource{
						{"Salle Mistral", "Salle de réunion", "available"},
						{"Salle Rhône", "Salle de réunion", "unavailable"},
						{"Bureau Lyon 01", "Bureau individuel", "available"},
						{"Espace Cowork Lyon", "Espace coworking", "available"},
						{"PC Lyon 01", "PC", "available"},
						{"PC Lyon 02", "PC", "available"},
						{"Projecteur Lyon 01", "Projecteur", "available"},
						{"Imprimante Lyon 01", "Imprimante", "available"},
					},
				},
			},
		},
		{
			name: "InnoSpace Agency",
			admin: seedUser{"admin.innospace@spacebook.fr", "admin_innospace", 600200300, "admin"},
			users: []seedUser{
				{"sophie.bernard@innospace.fr", "sophie_bernard", 600222001, "user"},
				{"thomas.petit@innospace.fr", "thomas_petit", 600222002, "user"},
				{"emma.richard@innospace.fr", "emma_richard", 600222003, "user"},
			},
			places: []seedPlace{
				{
					name:        "Campus Marseille",
					description: "Campus innovation sur le Vieux-Port",
					resources: []seedResource{
						{"Lab Innovation", "Laboratoire", "available"},
						{"Salle Méditerranée", "Salle de réunion", "available"},
						{"Salle Calanques", "Salle de conférence", "available"},
						{"Open Space Sud", "Espace coworking", "available"},
						{"Bureau Directeur", "Bureau individuel", "available"},
						{"PC Marseille 01", "PC", "available"},
						{"PC Marseille 02", "PC", "available"},
						{"PC Marseille 03", "PC", "available"},
						{"Projecteur Marseille 01", "Projecteur", "available"},
						{"Projecteur Marseille 02", "Projecteur", "unavailable"},
						{"Imprimante Marseille 01", "Imprimante", "available"},
					},
				},
				{
					name:        "Espace Bordeaux",
					description: "Espace de travail au cœur des Chartrons",
					resources: []seedResource{
						{"Salle Gironde", "Salle de réunion", "available"},
						{"Bureau B01", "Bureau individuel", "available"},
						{"Bureau B02", "Bureau individuel", "available"},
						{"Cowork Bordeaux", "Espace coworking", "available"},
						{"PC Bordeaux 01", "PC", "available"},
						{"PC Bordeaux 02", "PC", "available"},
						{"Projecteur Bordeaux 01", "Projecteur", "available"},
						{"Imprimante Bordeaux 01", "Imprimante", "available"},
						{"Imprimante Bordeaux 02", "Imprimante", "unavailable"},
					},
				},
			},
		},
		{
			name: "Digital Hub",
			admin: seedUser{"admin.digitalhub@spacebook.fr", "admin_digitalhub", 600300400, "admin"},
			users: []seedUser{
				{"lucas.simon@digitalhub.fr", "lucas_simon", 600333001, "user"},
				{"jade.michel@digitalhub.fr", "jade_michel", 600333002, "user"},
				{"noah.garcia@digitalhub.fr", "noah_garcia", 600333003, "user"},
				{"lea.martinez@digitalhub.fr", "lea_martinez", 600333004, "user"},
				{"hugo.sanchez@digitalhub.fr", "hugo_sanchez", 600333005, "user"},
			},
			places: []seedPlace{
				{
					name:        "Hub Nantes – Île de Nantes",
					description: "Hub numérique au cœur de l'île de Nantes",
					resources: []seedResource{
						{"Salle Loire", "Salle de réunion", "available"},
						{"Salle Atlantique", "Salle de réunion", "available"},
						{"Auditorium Digital", "Salle de conférence", "available"},
						{"Lab Dev", "Laboratoire", "available"},
						{"Open Desk 1", "Espace coworking", "available"},
						{"Open Desk 2", "Espace coworking", "available"},
						{"Bureau Senior 1", "Bureau individuel", "available"},
						{"PC Nantes 01", "PC", "available"},
						{"PC Nantes 02", "PC", "available"},
						{"PC Nantes 03", "PC", "available"},
						{"PC Nantes 04", "PC", "unavailable"},
						{"Projecteur Nantes 01", "Projecteur", "available"},
						{"Projecteur Nantes 02", "Projecteur", "available"},
						{"Imprimante Nantes 01", "Imprimante", "available"},
						{"Imprimante Nantes 02", "Imprimante", "available"},
					},
				},
				{
					name:        "Antenne Rennes",
					description: "Bureau satellite à Rennes",
					resources: []seedResource{
						{"Salle Bretagne", "Salle de réunion", "available"},
						{"Bureau Rennes 01", "Bureau individuel", "available"},
						{"Cowork Rennes", "Espace coworking", "available"},
						{"PC Rennes 01", "PC", "available"},
						{"PC Rennes 02", "PC", "available"},
						{"Projecteur Rennes 01", "Projecteur", "available"},
						{"Imprimante Rennes 01", "Imprimante", "available"},
					},
				},
			},
		},
	}

	for _, sc := range companies {
		// Company
		var company models.Company
		if db.Where("name = ?", sc.name).First(&company).Error != nil {
			company = models.Company{ID: uuid.New(), Name: sc.name}
			db.Create(&company)
			fmt.Printf("\n🏢 Entreprise créée : %s\n", sc.name)
		} else {
			fmt.Printf("\n⏭️  Entreprise existante : %s\n", sc.name)
		}

		// Admin
		var adminUser models.User
		if db.Where("email = ?", sc.admin.email).First(&adminUser).Error != nil {
			adminUser = models.User{
				ID:       uuid.New(),
				Email:    sc.admin.email,
				Username: sc.admin.username,
				Phone:    sc.admin.phone,
				Password: hash("Admin1234!"),
				Role:     "admin",
				Status:   "Actif",
			}
			db.Create(&adminUser)
			fmt.Printf("  👤 Admin créé : %s\n", sc.admin.email)
		} else {
			fmt.Printf("  ⏭️  Admin existant : %s\n", sc.admin.email)
		}

		// Places
		for i, sp := range sc.places {
			var place models.Place
			if db.Where("name = ? AND company_id = ?", sp.name, company.ID).First(&place).Error != nil {
				place = models.Place{
					ID:          uuid.New(),
					CompanyID:   company.ID,
					Name:        sp.name,
					Description: sp.description,
				}
				db.Create(&place)
				fmt.Printf("  📍 Lieu créé : %s\n", sp.name)
			} else {
				fmt.Printf("  ⏭️  Lieu existant : %s\n", sp.name)
			}

			// Link admin to first place only
			if i == 0 {
				var pu models.PlaceUser
				if db.Where("place_id = ? AND user_id = ?", place.ID, adminUser.ID).First(&pu).Error != nil {
					db.Create(&models.PlaceUser{ID: uuid.New(), PlaceID: place.ID, UserID: adminUser.ID})
				}
			}

			// Resources
			for _, sr := range sp.resources {
				var resource models.Resource
				if db.Where("name = ? AND place_id = ?", sr.name, place.ID).First(&resource).Error != nil {
					catID := catMap[sr.category]
					resource = models.Resource{
						ID:         uuid.New(),
						PlaceID:    &place.ID,
						CategoryID: &catID,
						Name:       sr.name,
						Status:     sr.status,
					}
					db.Create(&resource)
					fmt.Printf("    🖥️  Ressource créée : %s\n", sr.name)
				}
			}
		}

		// Regular users — linked to first place
		var firstPlace models.Place
		db.Where("company_id = ?", company.ID).Order("created_at ASC").First(&firstPlace)

		for _, su := range sc.users {
			var user models.User
			if db.Where("email = ?", su.email).First(&user).Error != nil {
				user = models.User{
					ID:       uuid.New(),
					Email:    su.email,
					Username: su.username,
					Phone:    su.phone,
					Password: hash("User1234!"),
					Role:     "user",
					Status:   "Actif",
				}
				db.Create(&user)
				db.Create(&models.PlaceUser{ID: uuid.New(), PlaceID: firstPlace.ID, UserID: user.ID})
				fmt.Printf("  👤 Utilisateur créé : %s\n", su.email)
			} else {
				fmt.Printf("  ⏭️  Utilisateur existant : %s\n", su.email)
			}
		}
	}

	// =====================
	// RESERVATIONS D'EXEMPLE
	// =====================
	fmt.Println("\n📅 Création des réservations d'exemple...")

	var alice models.User
	var neptune models.Resource
	if db.Where("email = ?", "alice.martin@techcorp.fr").First(&alice).Error == nil &&
		db.Where("name = ?", "Salle Neptune").First(&neptune).Error == nil {

		slots := []struct{ start, end time.Time }{
			{time.Now().Add(24 * time.Hour).Truncate(time.Hour), time.Now().Add(25 * time.Hour).Truncate(time.Hour)},
			{time.Now().Add(48 * time.Hour).Truncate(time.Hour), time.Now().Add(50 * time.Hour).Truncate(time.Hour)},
			{time.Now().Add(-48 * time.Hour).Truncate(time.Hour), time.Now().Add(-47 * time.Hour).Truncate(time.Hour)},
		}
		statuses := []string{"pending", "approved", "rejected"}

		for i, slot := range slots {
			var count int64
			db.Model(&models.Reservation{}).
				Where("user_id = ? AND ressource_id = ? AND start_datetime = ?", alice.ID, neptune.ID, slot.start).
				Count(&count)
			if count == 0 {
				resID := neptune.ID
				db.Create(&models.Reservation{
					ID:            uuid.New(),
					UserID:        alice.ID,
					RessourceID:   &resID,
					StartDatetime: slot.start,
					EndDatetime:   slot.end,
					Status:        statuses[i],
				})
				fmt.Printf("  📅 Réservation créée (%s) pour alice — %s\n", statuses[i], slot.start.Format("02/01 15h04"))
			}
		}
	}

	// =====================
	// NOTIFICATIONS D'EXEMPLE
	// =====================
	fmt.Println("\n🔔 Création des notifications d'exemple...")

	var adminTC models.User
	if db.Where("email = ?", "admin.techcorp@spacebook.fr").First(&adminTC).Error == nil {
		notifs := []struct{ msg, typ string }{
			{"alice.martin@techcorp.fr a demandé à rejoindre votre espace", "user_registration"},
			{"Nouvelle demande de réservation de alice_martin pour Salle Neptune", "reservation"},
		}
		for _, n := range notifs {
			var count int64
			db.Model(&models.Notification{}).Where("message = ? AND user_id = ?", n.msg, adminTC.ID).Count(&count)
			if count == 0 {
				adminID := adminTC.ID
				db.Create(&models.Notification{UserID: &adminID, Type: n.typ, Message: n.msg, IsRead: false})
				fmt.Printf("  🔔 Notification créée pour admin_techcorp\n")
			}
		}
	}

	fmt.Println("\n✅ Seed terminé !")
	printSummary()
}

func printSummary() {
	db := config.DB
	var cc, cp, cu, cr, ccat, cpu int64
	db.Model(&models.Company{}).Count(&cc)
	db.Model(&models.Place{}).Count(&cp)
	db.Model(&models.User{}).Count(&cu)
	db.Model(&models.Resource{}).Count(&cr)
	db.Model(&models.Category{}).Count(&ccat)
	db.Model(&models.PlaceUser{}).Count(&cpu)

	fmt.Println("\n📊 État de la base de données :")
	fmt.Printf("   Entreprises : %d\n", cc)
	fmt.Printf("   Lieux       : %d\n", cp)
	fmt.Printf("   Utilisateurs: %d\n", cu)
	fmt.Printf("   Ressources  : %d\n", cr)
	fmt.Printf("   Catégories  : %d\n", ccat)
	fmt.Printf("   PlaceUsers  : %d\n", cpu)
}
