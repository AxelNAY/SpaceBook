import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import logo from "../assets/SpaceBook_logo.png";
import LoginModal from "./LoginModal";
import RegisterModal from "./RegisterModal";

export default function Navbar() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const [showUserMenu, setShowUserMenu] = useState(false);
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    setShowUserMenu(false);
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

  return (
    <>
      <header className="header">
        <Link to={isAdmin ? "/admin/resources/create" : "/"} className="header-logo">
          <img src={logo} alt="SpaceBook" />
        </Link>

        <nav className="header-nav">
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
                        Mes Reservations
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
                        Reservations
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
                        Creation
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
                    Deconnexion
                  </button>
                </div>
              )}
            </div>
          )}
        </nav>
      </header>

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
