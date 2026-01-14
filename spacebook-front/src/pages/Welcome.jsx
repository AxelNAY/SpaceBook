import { useState } from "react";
import RegisterModal from "../components/RegisterModal";
import LoginModal from "../components/LoginModal";

export default function Welcome() {
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const [showLoginModal, setShowLoginModal] = useState(false);

  const switchToLogin = () => {
    setShowRegisterModal(false);
    setShowLoginModal(true);
  };

  const switchToRegister = () => {
    setShowLoginModal(false);
    setShowRegisterModal(true);
  };

  return (
    <>
      <div className="welcome-container">
        <h1 className="welcome-title">Bienvenue sur SpaceBook !</h1>

        <p className="welcome-text">
          En tant qu'utilisateur, vous pouvez réserver des salles et du matériel.
          <br />
          En tant qu'administrateur, vous pouvez gérer les demandes des utilisateurs
          et administrer les salles et le matériel.
        </p>

        <button
          className="btn btn-primary"
          onClick={() => setShowRegisterModal(true)}
        >
          Inscription
        </button>
      </div>

      {showRegisterModal && (
        <RegisterModal
          onClose={() => setShowRegisterModal(false)}
          onSwitchToLogin={switchToLogin}
        />
      )}

      {showLoginModal && (
        <LoginModal
          onClose={() => setShowLoginModal(false)}
          onSwitchToRegister={switchToRegister}
        />
      )}
    </>
  );
}
