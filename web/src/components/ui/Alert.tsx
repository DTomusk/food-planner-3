import { CircleCheckBig, Info, OctagonX, TriangleAlert, X } from "lucide-react";
import type { JSX } from "react/jsx-dev-runtime";
import Inline from "../layout/Inline";
import { useEffect, useState } from "react";

type AlertType = 'info' | 'error' | 'success' | 'warning';

interface AlertProps {
    message: string;
    type?: AlertType;
    duration?: number;
    onClose?: () => void;
    closable?: boolean;
}

const alertStyles: Record<AlertType, string> = {
    error: 'text-red-600 bg-red-50 border-red-200',
    success: 'text-green-600 bg-green-50 border-green-200',
    warning: 'text-yellow-600 bg-yellow-50 border-yellow-200',
    info: 'text-blue-600 bg-blue-50 border-blue-200'
};

const alertPrefixes: Record<AlertType, JSX.Element> = {
    error: <OctagonX/>,
    success: <CircleCheckBig/>,
    warning: <TriangleAlert/>,
    info: <Info/>
};

export default function Alert({ message, type = 'error', duration, onClose, closable = false }: AlertProps) {
    const [visible, setVisible] = useState(true);

    const handleClose = () => {
        setVisible(false);
        onClose?.();
    };

    useEffect(() => {
        if (!duration) return;
        const timer = setTimeout(() => {
            setVisible(false);
            onClose?.();
        }, duration);
        return () => clearTimeout(timer);
    }, [duration, onClose]);

    if (!visible) return null;

    const styleClass = alertStyles[type];
    const prefix = alertPrefixes[type];

    return (
        <div className={`font-bold p-3 rounded-md border ${styleClass}`}>
            <div className="flex justify-between items-start align-center">
            <Inline>
                {prefix}
                <span>{message}</span>
            </Inline>
            {closable && (
                <button
                    type="button"
                    className="ml-4 text-current hover:text-opacity-70 cursor-pointer"
                    onClick={handleClose}
                    aria-label="Close alert"
                >
                <X size={20} />
            </button>
            )}
            </div>
        </div>
    );
}