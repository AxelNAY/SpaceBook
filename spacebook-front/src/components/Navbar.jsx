import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import logo from "../assets/SpaceBook_logo.png";
import LoginModal from "./LoginModal";
import RegisterModal from "./RegisterModal";

export default function Navbar() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const [showUserMenu, setShowUserMenu] = useState(false);
  const [showMobileMenu, setShowMobileMenu] = useState(false);
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    setShowUserMenu(false);
    setShowMobileMenu(false);
    navigate("/");
  };

  const switchToRegister = () => {
    setShowLoginModal(false);
    setShowRegisterModal(true);
  };

  const switchToLogin = () => {
    setShowRegisterModal(false);
    setShowLoginModal(true);
  };

  const closeMobileMenu = () => setShowMobileMenu(false);

  return (
    <>
      <header className="header">
        <Link
          to={!isAuthenticated ? "/" : isAdmin ? "/admin/reservations" : "/ressources"}
          className="header-logo"
        >
          <img src={logo} alt="SpaceBook" />
        </Link>

        <nav className="header-nav header-nav-desktop">
          {isAuthenticated && (
            <Link to="/notifications" className="header-nav-item">
              Notifications
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
              </svg>
            </Link>
          )}

          {!isAuthenticated ? (
            <button
              className="header-nav-item"
              onClick={() => setShowLoginModal(true)}
              style={{ background: "none", border: "none" }}
            >
              Connexion
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="8" r="4" />
                <path d="M20 21a8 8 0 1 0-16 0" />
              </svg>
            </button>
          ) : (
            <div className="user-menu">
              <button
                className="header-nav-item"
                onClick={() => setShowUserMenu(!showUserMenu)}
                style={{ background: "none", border: "none" }}
              >
                {isAdmin ? "Admin" : user?.username || "Utilisateur"}
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="12" cy="8" r="4" />
                  <path d="M20 21a8 8 0 1 0-16 0" />
                </svg>
              </button>

              {showUserMenu && (
                <div className="user-menu-dropdown">
                  {!isAdmin && (
                    <>
                      <Link
                        to="/ressources"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Ressources
                      </Link>
                      <Link
                        to="/mes-reservations"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Mes Réservations
                      </Link>
                    </>
                  )}
                  {isAdmin && (
                    <>
                      <Link
                        to="/admin/reservations"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Réservations
                      </Link>
                      <Link
                        to="/admin/resources"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Ressources
                      </Link>
                      <Link
                        to="/admin/resources/create"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Création
                      </Link>
                      <Link
                        to="/admin/places"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Lieux
                      </Link>
                      <Link
                        to="/admin/users"
                        className="user-menu-item btn btn-primary"
                        onClick={() => setShowUserMenu(false)}
                      >
                        Utilisateurs
                      </Link>
                    </>
                  )}
                  <button
                    className="user-menu-item btn btn-danger"
                    onClick={handleLogout}
                  >
                    Déconnexion
                  </button>
                </div>
              )}
            </div>
          )}
        </nav>

        <button
          className="hamburger-btn"
          onClick={() => setShowMobileMenu(!showMobileMenu)}
          aria-label="Menu"
        >
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
            {showMobileMenu ? (
              <>
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </>
            ) : (
              <>
                <line x1="3" y1="6" x2="21" y2="6" />
                <line x1="3" y1="12" x2="21" y2="12" />
                <line x1="3" y1="18" x2="21" y2="18" />
              </>
            )}
          </svg>
        </button>
      </header>

      {showMobileMenu && (
        <nav className="mobile-menu">
          {isAuthenticated && (
            <Link to="/notifications" className="mobile-menu-item" onClick={closeMobileMenu}>
              Notifications
            </Link>
          )}
          {!isAuthenticated ? (
            <button
              className="mobile-menu-item"
              onClick={() => { setShowLoginModal(true); closeMobileMenu(); }}
            >
              Connexion
            </button>
          ) : (
            <>
              {!isAdmin && (
                <>
                  <Link to="/ressources" className="mobile-menu-item" onClick={closeMobileMenu}>Ressources</Link>
                  <Link to="/mes-reservations" className="mobile-menu-item" onClick={closeMobileMenu}>Mes Réservations</Link>
                </>
              )}
              {isAdmin && (
                <>
                  <Link to="/admin/reservations" className="mobile-menu-item" onClick={closeMobileMenu}>Réservations</Link>
                  <Link to="/admin/resources" className="mobile-menu-item" onClick={closeMobileMenu}>Ressources</Link>
                  <Link to="/admin/resources/create" className="mobile-menu-item" onClick={closeMobileMenu}>Création</Link>
                  <Link to="/admin/places" className="mobile-menu-item" onClick={closeMobileMenu}>Lieux</Link>
                  <Link to="/admin/users" className="mobile-menu-item" onClick={closeMobileMenu}>Utilisateurs</Link>
                </>
              )}
              <button className="mobile-menu-item mobile-menu-logout" onClick={handleLogout}>
                Déconnexion
              </button>
            </>
          )}
        </nav>
      )}

      {showLoginModal && (
        <LoginModal
          onClose={() => setShowLoginModal(false)}
          onSwitchToRegister={switchToRegister}
        />
      )}

      {showRegisterModal && (
        <RegisterModal
          onClose={() => setShowRegisterModal(false)}
          onSwitchToLogin={switchToLogin}
        />
      )}
    </>
  );
}
