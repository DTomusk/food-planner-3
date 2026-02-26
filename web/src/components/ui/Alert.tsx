import { CircleCheckBig, Info, OctagonX, TriangleAlert } from "lucide-react";
import type { JSX } from "react/jsx-dev-runtime";
import Inline from "../layout/Inline";

type AlertType = 'info' | 'error' | 'success' | 'warning';

interface AlertProps {
    message: string;
    type?: AlertType;
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

export default function Alert({ message, type = 'error' }: AlertProps) {
    const styleClass = alertStyles[type];
    const prefix = alertPrefixes[type];
    
    return (
        <div className={`font-bold p-3 rounded-md border ${styleClass}`}>
            <Inline>
                {prefix}
                <span>{message}</span>
            </Inline>
        </div>
    );
}