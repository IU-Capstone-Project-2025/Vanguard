import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { QRCodeCanvas } from 'qrcode.react';
import styles from './QRCodeModal.module.css';

const QRCodeModal = ({ code, onClose }) => {
  return (
    <AnimatePresence>
      <motion.div
        className={styles['modal-overlay']}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        onClick={onClose}
      >
        <motion.div
          className={styles['modal-wrapper']}
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          exit={{ scale: 0 }}
          onClick={(e) => e.stopPropagation()}
        >
          <div className={styles['modal-navbar']}>
            <button onClick={onClose}>⤫</button>
          </div>
          <div className={styles['qr-container']}>
            <QRCodeCanvas
              value={code}
              size={620}
              bgColor="#ffffff"
              fgColor="#000000"
              level="H"
              includeMargin={true}
            />
            <p>Scan to join session</p>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
};

export default QRCodeModal;
