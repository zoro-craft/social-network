'use client';

import { apiRequest } from "@/utils/api";
import { useState, useEffect } from "react";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import styles from '../styles/user.module.css';

export default function userlistpage() {
    const { auth, loading } = useRequireAuth();
    const [users, setUsers] = useState([])
    useEffect(() => {
        if (!auth.authenticated) return
        const Fetchusers = async () => {
            try {
                const res = await apiRequest("/users/search")
                setUsers(res || [])
            } catch {
                setUsers([])
            }
        }
        Fetchusers();
    }, [auth.authenticated])

    if (loading) {
        return (
            <div className="glass-feed-page">
                <p className={styles.mutedText}>Loading ...</p>
            </div>
        );
    }

    return  (
         <></>
    )
}