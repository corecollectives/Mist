import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FullScreenLoading } from '@/shared/components';
import { useUsersStore } from './store';
import { UserCard, CreateUserModal } from './components';
import { canManageUsers } from './utils';
import type { CreateUserData, User } from '@/types';
import { useAuth } from '@/context/AuthContext';

export default function UsersPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { user } = useAuth()
  const {
    users,
    fetchUsers,
    createUser,
    isLoading,
    error,
    clearError
  } = useUsersStore();

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  useEffect(() => {
    if (error) {
      toast.error(error);
      clearError();
    }
  }, [error, clearError]);

  const handleCreateUser = async (userData: CreateUserData) => {
    const success = await createUser(userData);
    if (success) {
      toast.success('User created successfully');
      setIsModalOpen(false);
    } else {
    }
  };
  const handleUserClick = (selectedUser: User) => {
    console.log('User clicked:', selectedUser);
  };

  if (isLoading && users.length === 0) {
    return <FullScreenLoading />;
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="flex items-center justify-between py-6 border-b border-border shrink-0">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Users
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage users and their permissions
          </p>
        </div>
        <Button
          onClick={() => setIsModalOpen(true)}
          disabled={!canManageUsers(user)}
          className="transition-colors"
        >
          Add User
        </Button>
      </div>

      {/* Content */}
      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {users.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12">
          <p className="text-muted-foreground text-lg mb-4">No users found</p>
          {canManageUsers(user) && (
            <Button onClick={() => setIsModalOpen(true)}>
              Create First User
            </Button>
          )}
        </div>
      ) : (
        <div className="grid gap-4 py-6 md:grid-cols-2 lg:grid-cols-3">
          {users.map((userItem: User) => (
            <UserCard
              key={userItem.id}
              user={userItem}
              onClick={handleUserClick}
            />
          ))}
        </div>
      )}

      {/* Create User Modal */}
      <CreateUserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateUser}
        currentUser={user}
      />
    </div>
  );
}
